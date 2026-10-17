package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"bluespot/internal/config"
	"bluespot/internal/model"
	"bluespot/internal/repository"
)

var ErrInvalidAttachFile = errors.New("invalid attach file")
var ErrInvalidConfigFileSecret = errors.New("invalid config file upload secret")

type AttachService struct {
	repo repository.AttachRepository
}

func NewAttachService(repo repository.AttachRepository) *AttachService {
	return &AttachService{repo: repo}
}

func (s *AttachService) Upload(file *multipart.FileHeader) (*model.AttachUploadResponse, error) {
	if file == nil {
		return nil, ErrInvalidAttachFile
	}

	originName := strings.TrimSpace(filepath.Base(file.Filename))
	if originName == "" || originName == "." {
		return nil, ErrInvalidAttachFile
	}

	fileName, err := generateAttachFileName(originName)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(file, fileName); err != nil {
		return nil, err
	}

	return &model.AttachUploadResponse{
		FileName:   fileName,
		OriginName: originName,
		URL:        buildAttachURL(fileName),
	}, nil
}

func generateAttachFileName(originName string) (string, error) {
	ext := filepath.Ext(originName)
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("generate random suffix failed: %w", err)
	}
	return fmt.Sprintf("%d%06d%s", time.Now().UnixNano(), n.Int64(), ext), nil
}

func buildAttachURL(fileName string) string {
	baseURL := strings.TrimRight(strings.TrimSpace(config.GetAttachViewBaseURL()), "/")
	if baseURL == "" {
		return url.PathEscape(fileName)
	}
	return baseURL + "/" + url.PathEscape(fileName)
}

// UploadConfigFile 校验上传密钥后按原始文件名覆盖配置文件。
func (s *AttachService) UploadConfigFile(file *multipart.FileHeader, secret string) (*model.AttachUploadResponse, error) {
	configuredSecret := config.GetAttachConfigFileUploadSecret()
	if configuredSecret == "" || secret == "" || secret != configuredSecret {
		return nil, ErrInvalidConfigFileSecret
	}
	if file == nil {
		return nil, ErrInvalidAttachFile
	}

	// 同时处理 Unix 和 Windows 路径，文件始终保存在配置目录内。
	fileName := filepath.Base(strings.ReplaceAll(file.Filename, "\\", "/"))
	if strings.TrimSpace(fileName) == "" || fileName == "." || fileName == ".." || strings.ContainsAny(fileName, "/\x00") {
		return nil, ErrInvalidAttachFile
	}
	if err := s.repo.SaveConfigFile(file, fileName); err != nil {
		return nil, err
	}

	baseURL := strings.TrimRight(strings.TrimSpace(config.GetAttachConfigFileViewBaseURL()), "/")
	fileURL := url.PathEscape(fileName)
	if baseURL != "" {
		fileURL = baseURL + "/" + fileURL
	}
	return &model.AttachUploadResponse{
		FileName:   fileName,
		OriginName: fileName,
		URL:        fileURL,
	}, nil
}
