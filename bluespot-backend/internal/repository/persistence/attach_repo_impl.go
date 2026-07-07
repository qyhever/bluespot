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
