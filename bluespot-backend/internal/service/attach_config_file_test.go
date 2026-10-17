package service

import (
	"errors"
	"mime/multipart"
	"testing"

	"bluespot/internal/config"
)

func TestUploadConfigFileRejectsInvalidInputBeforeSaving(t *testing.T) {
	previous := config.GlobalConfig
	t.Cleanup(func() { config.GlobalConfig = previous })
	config.GlobalConfig = &config.Config{Attach: config.AttachConfig{UploadConfigFileSecret: "secret"}}
	// 未提供 repository，校验失败必须在调用存储层之前返回。
	svc := NewAttachService(nil)
	for _, name := range []string{"", " ", ".", "..", "/", `\`, "../..", "bad\x00.json"} {
		_, err := svc.UploadConfigFile(&multipart.FileHeader{Filename: name}, "secret")
		if !errors.Is(err, ErrInvalidAttachFile) {
			t.Fatalf("name %q: error=%v", name, err)
		}
	}
	if _, err := svc.UploadConfigFile(nil, "secret"); !errors.Is(err, ErrInvalidAttachFile) {
		t.Fatalf("nil file: %v", err)
	}
	config.GlobalConfig = nil
	if _, err := svc.UploadConfigFile(&multipart.FileHeader{Filename: "config.json"}, "secret"); !errors.Is(err, ErrInvalidConfigFileSecret) {
		t.Fatalf("nil config: %v", err)
	}
}
