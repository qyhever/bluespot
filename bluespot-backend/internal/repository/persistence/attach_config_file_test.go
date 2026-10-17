package persistence

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"bluespot/internal/config"
)

func TestSaveConfigFileFailurePreservesDestination(t *testing.T) {
	previous := config.GlobalConfig
	t.Cleanup(func() { config.GlobalConfig = previous })
	dir := t.TempDir()
	config.GlobalConfig = &config.Config{Attach: config.AttachConfig{UploadConfigFileDirPath: dir}}
	repo := NewAttachRepository()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("new content")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	payload := append([]byte(nil), body.Bytes()...)
	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 强制 multipart 使用磁盘文件，然后移除源文件以模拟读取失败。
	if err := req.ParseMultipartForm(0); err != nil {
		t.Fatal(err)
	}
	header := req.MultipartForm.File["file"][0]
	if err := req.MultipartForm.RemoveAll(); err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(dir, "config.json")
	if err := os.WriteFile(destination, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveConfigFile(header, "config.json"); err == nil {
		t.Fatal("expected source read failure")
	}
	content, err := os.ReadFile(destination)
	if err != nil || string(content) != "old content" {
		t.Fatalf("old file changed: %q, %v", content, err)
	}
	// 有效文件写入后无法替换已有目录，临时文件也应清理。
	memoryReq := httptest.NewRequest("POST", "/", bytes.NewReader(payload))
	memoryReq.Header.Set("Content-Type", writer.FormDataContentType())
	if err := memoryReq.ParseMultipartForm(1024); err != nil {
		t.Fatal(err)
	}
	defer memoryReq.MultipartForm.RemoveAll()
	validHeader := memoryReq.MultipartForm.File["file"][0]
	if err := os.Mkdir(filepath.Join(dir, "blocked"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveConfigFile(validHeader, "blocked"); err == nil {
		t.Fatal("expected rename failure")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("temporary files leaked: %v", entries)
	}
	for _, name := range []string{"", ".", "..", "../escape", `..\escape`} {
		if err := repo.SaveConfigFile(&multipart.FileHeader{}, name); err == nil {
			t.Fatalf("accepted invalid name %q", name)
		}
	}
}
