package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"bluespot/internal/config"
	"bluespot/internal/model"
	"bluespot/internal/repository"
)

var (
	ErrInvalidUploadParam = errors.New("invalid upload param")
	ErrUploadIDMismatch   = errors.New("upload id mismatch")
	ErrChunkMissing       = errors.New("chunk missing")
)

type UploadService struct {
	repo repository.UploadRepository
}

func NewUploadService(repo repository.UploadRepository) *UploadService {
	return &UploadService{repo: repo}
}

func (s *UploadService) Verify(userID uint64, req model.UploadVerifyRequest) (*model.UploadVerifyResponse, error) {
	fileMD5, fileName, err := normalizeVerifyRequest(req)
	if err != nil {
		return nil, err
	}

	uploadID := GenerateUploadID(fileMD5, userID, config.GetAttachChunkDirSalt())
	finalFileName := buildLargeFileName(fileMD5, fileName)
	url := buildLargeFileURL(finalFileName)

	exists, err := s.repo.FinalFileExists(finalFileName)
	if err != nil {
		return nil, err
	}
	if exists {
		return &model.UploadVerifyResponse{
			IsExists:       true,
			URL:            url,
			UploadID:       uploadID,
			UploadedChunks: []int{},
		}, nil
	}

	if err := s.repo.SaveMeta(uploadID, repository.ChunkUploadMeta{FileName: fileName}); err != nil {
		return nil, err
	}
	chunks, err := s.repo.ListUploadedChunks(uploadID)
	if err != nil {
		return nil, err
	}
	return &model.UploadVerifyResponse{
		IsExists:       false,
		URL:            "",
		UploadID:       uploadID,
		UploadedChunks: chunks,
	}, nil
}

func (s *UploadService) SaveChunk(userID uint64, uploadID, fileMD5, fileName string, chunkIndex int, chunk *multipart.FileHeader) error {
	fileMD5 = strings.TrimSpace(fileMD5)
	fileName = safeBaseName(fileName)
	if !isValidMD5Hex(fileMD5) || fileName == "" || chunkIndex < 0 || chunk == nil {
		return ErrInvalidUploadParam
	}
	if !s.isExpectedUploadID(userID, fileMD5, uploadID) {
		return ErrUploadIDMismatch
	}
	if err := s.repo.SaveMeta(uploadID, repository.ChunkUploadMeta{FileName: fileName}); err != nil {
		return err
	}
	return s.repo.SaveChunk(uploadID, chunkIndex, chunk)
}

func (s *UploadService) Merge(userID uint64, req model.UploadMergeRequest) (*model.UploadMergeResponse, error) {
	fileMD5 := strings.TrimSpace(req.FileMD5)
	if !isValidMD5Hex(fileMD5) || req.ChunkLength <= 0 {
		return nil, ErrInvalidUploadParam
	}
	if !s.isExpectedUploadID(userID, fileMD5, req.UploadID) {
		return nil, ErrUploadIDMismatch
	}

	existingFinalFileName, err := s.repo.FindFinalFileNameByMD5(fileMD5)
	if err != nil {
		return nil, err
	}
	if existingFinalFileName != "" {
		return &model.UploadMergeResponse{
			URL: buildLargeFileURL(existingFinalFileName),
			Msg: "合并成功",
		}, nil
	}

	meta, err := s.repo.ReadMeta(req.UploadID)
	if err != nil {
		return nil, err
	}
	finalFileName := buildLargeFileName(fileMD5, meta.FileName)
	finalURL := buildLargeFileURL(finalFileName)

	exists, err := s.repo.FinalFileExists(finalFileName)
	if err != nil {
		return nil, err
	}
	if exists {
		return &model.UploadMergeResponse{
			URL: finalURL,
			Msg: "合并成功",
		}, nil
	}

	for i := 0; i < req.ChunkLength; i++ {
		exists, err := s.repo.ChunkExists(req.UploadID, i)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("%w: %d", ErrChunkMissing, i)
		}
	}
	if err := s.repo.MergeChunks(req.UploadID, req.ChunkLength, finalFileName); err != nil {
		return nil, err
	}
	if err := s.repo.CleanupChunks(req.UploadID); err != nil {
		return nil, err
	}

	return &model.UploadMergeResponse{
		URL: finalURL,
		Msg: "合并成功",
	}, nil
}

func GenerateUploadID(fileMD5 string, userID uint64, salt string) string {
	sum := md5.Sum([]byte(fmt.Sprintf("%s%d%s", fileMD5, userID, salt)))
	return hex.EncodeToString(sum[:])
}

func (s *UploadService) isExpectedUploadID(userID uint64, fileMD5, uploadID string) bool {
	expected := GenerateUploadID(fileMD5, userID, config.GetAttachChunkDirSalt())
	return strings.TrimSpace(uploadID) == expected
}

func normalizeVerifyRequest(req model.UploadVerifyRequest) (string, string, error) {
	fileMD5 := strings.TrimSpace(req.FileMD5)
	fileName := safeBaseName(req.FileName)
	if !isValidMD5Hex(fileMD5) || fileName == "" || req.FileSize <= 0 {
		return "", "", ErrInvalidUploadParam
	}
	return fileMD5, fileName, nil
}

func safeBaseName(fileName string) string {
	normalized := strings.ReplaceAll(fileName, "\\", "/")
	name := strings.TrimSpace(path.Base(normalized))
	if name == "" || name == "." || name == string(filepath.Separator) {
		return ""
	}
	return name
}

func isValidMD5Hex(value string) bool {
	if len(value) != 32 {
		return false
	}
	for _, r := range value {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
			continue
		}
		return false
	}
	return true
}

func buildLargeFileName(fileMD5, fileName string) string {
	return fileMD5 + filepath.Ext(fileName)
}

func buildLargeFileURL(fileName string) string {
	baseURL := strings.TrimRight(strings.TrimSpace(config.GetAttachLargeFileViewBaseURL()), "/")
	if baseURL == "" {
		return url.PathEscape(fileName)
	}
	return baseURL + "/" + url.PathEscape(fileName)
}
