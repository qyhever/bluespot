package persistence

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"bluespot/internal/config"
	"bluespot/internal/repository"
)

type AttachRepositoryImpl struct{}

func NewAttachRepository() repository.AttachRepository {
	return &AttachRepositoryImpl{}
}

func (r *AttachRepositoryImpl) Save(file *multipart.FileHeader, fileName string) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}
	fileName = strings.TrimSpace(fileName)
	if fileName == "" || fileName != filepath.Base(fileName) {
		return fmt.Errorf("invalid file name")
	}

	uploadDirPath := strings.TrimSpace(config.GetAttachUploadDirPath())
	if uploadDirPath == "" {
		return fmt.Errorf("upload dir path is empty")
	}
	if err := os.MkdirAll(uploadDirPath, 0755); err != nil {
		return fmt.Errorf("create upload dir failed: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("open upload file failed: %w", err)
	}
	defer src.Close()

	dstPath := filepath.Join(uploadDirPath, fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create upload file failed: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("save upload file failed: %w", err)
	}
	return nil
}

// SaveConfigFile 使用同目录临时文件和重命名，避免写入失败破坏旧文件。
func (r *AttachRepositoryImpl) SaveConfigFile(file *multipart.FileHeader, fileName string) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}
	if strings.TrimSpace(fileName) == "" || fileName == "." || fileName == ".." ||
		fileName != filepath.Base(fileName) || strings.ContainsAny(fileName, "/\\\x00") {
		return fmt.Errorf("invalid file name")
	}
	uploadDirPath := strings.TrimSpace(config.GetAttachConfigFileUploadDirPath())
	if uploadDirPath == "" {
		return fmt.Errorf("config file upload dir path is empty")
	}
	if err := os.MkdirAll(uploadDirPath, 0755); err != nil {
		return fmt.Errorf("create config file upload dir failed: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("open config upload file failed: %w", err)
	}
	defer src.Close()

	dst, err := os.CreateTemp(uploadDirPath, ".config-upload-*")
	if err != nil {
		return fmt.Errorf("create temporary config file failed: %w", err)
	}
	defer os.Remove(dst.Name())
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("write config file failed: %w", err)
	}
	if err := dst.Chmod(0644); err != nil {
		return fmt.Errorf("set config file permissions failed: %w", err)
	}
	if err := dst.Close(); err != nil {
		return fmt.Errorf("close config file failed: %w", err)
	}
	if err := os.Rename(dst.Name(), filepath.Join(uploadDirPath, fileName)); err != nil {
		return fmt.Errorf("replace config file failed: %w", err)
	}
	return nil
}
