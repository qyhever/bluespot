package router

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"bluespot/internal/config"
	"bluespot/internal/controller"
	"bluespot/internal/model"
)

func TestUploadConfigFile(t *testing.T) {
	previous := config.GlobalConfig
	t.Cleanup(func() { config.GlobalConfig = previous })
	dir := filepath.Join(t.TempDir(), "config")
	config.GlobalConfig = &config.Config{Attach: config.AttachConfig{
		UploadConfigFileDirPath: dir, UploadConfigFileSecret: " secret ",
		ViewConfigFileBaseURL: "https://example.test/config/",
	}}
	router := SetupRouter()
	upload := func(name, secret, content string) (controller.MyCode, string, model.AttachUploadResponse) {
		t.Helper()
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		if name != "" {
			part, err := writer.CreateFormFile("file", name)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := part.Write([]byte(content)); err != nil {
				t.Fatal(err)
			}
		}
		if secret != "" {
			if err := writer.WriteField("secret", secret); err != nil {
				t.Fatal(err)
			}
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/attach/uploadConfigFile", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		if response.Code != http.StatusOK {
			t.Fatalf("HTTP status = %d: %s", response.Code, response.Body.String())
		}
		var result struct {
			Code    controller.MyCode          `json:"code"`
			Message string                     `json:"message"`
			Data    model.AttachUploadResponse `json:"data"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		return result.Code, result.Message, result.Data
	}
	name := "配置 #1.json"
	code, _, data := upload(name, " secret ", "old")
	if code != controller.CodeSuccess {
		t.Fatalf("upload code = %d", code)
	}
	if data.FileName != name || data.OriginName != name || data.URL != "https://example.test/config/"+url.PathEscape(name) {
		t.Fatalf("unexpected response: %+v", data)
	}
	checkContent := func(want string) {
		t.Helper()
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil || string(content) != want {
			t.Fatalf("file = %q, err = %v; want %q", content, err, want)
		}
	}
	checkContent("old")
	for _, secret := range []string{"", "wrong", "secret"} {
		code, message, _ := upload(name, secret, "rejected")
		if code != controller.CodeInvalidParam || message != "上传密钥错误" {
			t.Fatalf("secret %q: code=%d message=%q", secret, code, message)
		}
		checkContent("old")
	}
	config.GlobalConfig.Attach.UploadConfigFileSecret = ""
	code, message, _ := upload(name, " secret ", "rejected")
	if code != controller.CodeInvalidParam || message != "上传密钥错误" {
		t.Fatalf("empty config secret: code=%d message=%q", code, message)
	}
	checkContent("old")
	config.GlobalConfig.Attach.UploadConfigFileSecret = " secret "
	for _, invalid := range []string{"", ".", ".."} {
		code, _, _ := upload(invalid, " secret ", "invalid")
		if code != controller.CodeInvalidParam {
			t.Fatalf("name %q: code=%d", invalid, code)
		}
	}
	code, _, _ = upload(name, " secret ", "new")
	if code != controller.CodeSuccess {
		t.Fatalf("overwrite code=%d", code)
	}
	checkContent("new")
	for _, pathName := range []string{"../../escape.json", `..\..\windows.json`} {
		code, _, data := upload(pathName, " secret ", "safe")
		if code != controller.CodeSuccess {
			t.Fatalf("path %q: code=%d", pathName, code)
		}
		content, err := os.ReadFile(filepath.Join(dir, data.FileName))
		if err != nil || string(content) != "safe" {
			t.Fatalf("path %q: content=%q err=%v", pathName, content, err)
		}
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(dir), "escape.json")); !os.IsNotExist(err) {
		t.Fatalf("file escaped upload dir: %v", err)
	}
	// 已有普通附件路由仍需 JWT。
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/attach/upload", nil))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("regular upload status=%d", response.Code)
	}
	// 配置目录指向普通文件时返回服务器错误。
	config.GlobalConfig.Attach.UploadConfigFileDirPath = filepath.Join(dir, name)
	code, _, _ = upload("other.json", " secret ", "content")
	if code != controller.CodeServerBusy {
		t.Fatalf("storage failure code=%d", code)
	}
	checkContent("new")
}
