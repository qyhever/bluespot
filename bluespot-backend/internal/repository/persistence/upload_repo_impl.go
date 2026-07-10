package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"bluespot/internal/config"
	"bluespot/internal/repository"
)

const uploadMetaFileName = ".meta.json"

type UploadRepositoryImpl struct{}

func NewUploadRepository() repository.UploadRepository {
	return &UploadRepositoryImpl{}
}

func (r *UploadRepositoryImpl) FinalFileExists(fileName string) (bool, error) {
	fileName, err := validateFinalFileName(fileName)
	if err != nil {
		return false, err
	}
	path, err := finalFilePath(fileName)
	if err != nil {
		return false, err
	}
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat final file failed: %w", err)
}

func (r *UploadRepositoryImpl) FindFinalFileNameByMD5(fileMD5 string) (string, error) {
	fileMD5 = strings.TrimSpace(fileMD5)
	if fileMD5 == "" || fileMD5 != filepath.Base(fileMD5) {
		return "", fmt.Errorf("invalid file md5")
	}
	root := strings.TrimSpace(config.GetAttachLargeFileUploadPath())
	if root == "" {
		return "", fmt.Errorf("large file upload path is empty")
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read large file dir failed: %w", err)
	}

	prefix := fileMD5 + "."
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name == fileMD5 || strings.HasPrefix(name, prefix) {
			return name, nil
		}
	}
	return "", nil
}

func (r *UploadRepositoryImpl) ListUploadedChunks(uploadID string) ([]int, error) {
	dir, err := chunkDir(uploadID)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []int{}, nil
		}
		return nil, fmt.Errorf("read chunk dir failed: %w", err)
	}

	chunks := make([]int, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		index, err := strconv.Atoi(entry.Name())
		if err != nil || index < 0 {
			continue
		}
		chunks = append(chunks, index)
	}
	sort.Ints(chunks)
	return chunks, nil
}

func (r *UploadRepositoryImpl) SaveChunk(uploadID string, chunkIndex int, file *multipart.FileHeader) error {
	if chunkIndex < 0 {
		return fmt.Errorf("invalid chunk index")
	}
	if file == nil {
		return fmt.Errorf("chunk file is nil")
	}
	dir, err := chunkDir(uploadID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create chunk dir failed: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("open chunk file failed: %w", err)
	}
	defer src.Close()

	dstPath := filepath.Join(dir, strconv.Itoa(chunkIndex))
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create chunk file failed: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("save chunk file failed: %w", err)
	}
	return nil
}

func (r *UploadRepositoryImpl) SaveMeta(uploadID string, meta repository.ChunkUploadMeta) error {
	if strings.TrimSpace(meta.FileName) == "" || meta.FileName != filepath.Base(meta.FileName) {
		return fmt.Errorf("invalid upload meta file name")
	}
	dir, err := chunkDir(uploadID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create chunk dir failed: %w", err)
	}
	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("marshal upload meta failed: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, uploadMetaFileName), data, 0644); err != nil {
		return fmt.Errorf("save upload meta failed: %w", err)
	}
	return nil
}

func (r *UploadRepositoryImpl) ReadMeta(uploadID string) (repository.ChunkUploadMeta, error) {
	dir, err := chunkDir(uploadID)
	if err != nil {
		return repository.ChunkUploadMeta{}, err
	}
	data, err := os.ReadFile(filepath.Join(dir, uploadMetaFileName))
	if err != nil {
		return repository.ChunkUploadMeta{}, fmt.Errorf("read upload meta failed: %w", err)
	}
	var meta repository.ChunkUploadMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return repository.ChunkUploadMeta{}, fmt.Errorf("unmarshal upload meta failed: %w", err)
	}
	if strings.TrimSpace(meta.FileName) == "" || meta.FileName != filepath.Base(meta.FileName) {
		return repository.ChunkUploadMeta{}, fmt.Errorf("invalid upload meta")
	}
	return meta, nil
}

func (r *UploadRepositoryImpl) ChunkExists(uploadID string, chunkIndex int) (bool, error) {
	if chunkIndex < 0 {
		return false, fmt.Errorf("invalid chunk index")
	}
	dir, err := chunkDir(uploadID)
	if err != nil {
		return false, err
	}
	info, err := os.Stat(filepath.Join(dir, strconv.Itoa(chunkIndex)))
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat chunk file failed: %w", err)
}

func (r *UploadRepositoryImpl) MergeChunks(uploadID string, chunkLength int, finalFileName string) error {
	if chunkLength <= 0 {
		return fmt.Errorf("invalid chunk length")
	}
	finalFileName, err := validateFinalFileName(finalFileName)
	if err != nil {
		return err
	}

	srcDir, err := chunkDir(uploadID)
	if err != nil {
		return err
	}
	dstPath, err := finalFilePath(finalFileName)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("create final file dir failed: %w", err)
	}

	tmpPath := dstPath + ".tmp"
	dst, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create final file failed: %w", err)
	}
	cleanupTmp := true
	defer func() {
		dst.Close()
		if cleanupTmp {
			_ = os.Remove(tmpPath)
		}
	}()

	for i := 0; i < chunkLength; i++ {
		chunkPath := filepath.Join(srcDir, strconv.Itoa(i))
		src, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("open chunk %d failed: %w", i, err)
		}
		if _, err := io.Copy(dst, src); err != nil {
			src.Close()
			return fmt.Errorf("merge chunk %d failed: %w", i, err)
		}
		if err := src.Close(); err != nil {
			return fmt.Errorf("close chunk %d failed: %w", i, err)
		}
	}
	if err := dst.Close(); err != nil {
		return fmt.Errorf("close final file failed: %w", err)
	}
	if err := os.Rename(tmpPath, dstPath); err != nil {
		return fmt.Errorf("rename final file failed: %w", err)
	}
	cleanupTmp = false
	return nil
}

func (r *UploadRepositoryImpl) DeleteFinalFile(fileName string) error {
	fileName, err := validateFinalFileName(fileName)
	if err != nil {
		return err
	}
	path, err := finalFilePath(fileName)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete final file failed: %w", err)
	}
	return nil
}

func (r *UploadRepositoryImpl) CleanupChunks(uploadID string) error {
	dir, err := chunkDir(uploadID)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("cleanup chunk dir failed: %w", err)
	}
	return nil
}

func chunkDir(uploadID string) (string, error) {
	uploadID = strings.TrimSpace(uploadID)
	if uploadID == "" || uploadID != filepath.Base(uploadID) {
		return "", fmt.Errorf("invalid upload id")
	}
	root := strings.TrimSpace(config.GetAttachChunkDirPath())
	if root == "" {
		return "", fmt.Errorf("chunk dir path is empty")
	}
	return filepath.Join(root, uploadID), nil
}

func validateFinalFileName(fileName string) (string, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" || fileName != filepath.Base(fileName) {
		return "", fmt.Errorf("invalid final file name")
	}
	return fileName, nil
}

func finalFilePath(fileName string) (string, error) {
	root := strings.TrimSpace(config.GetAttachLargeFileUploadPath())
	if root == "" {
		return "", fmt.Errorf("large file upload path is empty")
	}
	return filepath.Join(root, fileName), nil
}
