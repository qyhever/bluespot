package repository

import "mime/multipart"

type ChunkUploadMeta struct {
	FileName string
}

type UploadRepository interface {
	FinalFileExists(fileName string) (bool, error)
	FindFinalFileNameByMD5(fileMD5 string) (string, error)
	ListUploadedChunks(uploadID string) ([]int, error)
	SaveChunk(uploadID string, chunkIndex int, file *multipart.FileHeader) error
	SaveMeta(uploadID string, meta ChunkUploadMeta) error
	ReadMeta(uploadID string) (ChunkUploadMeta, error)
	ChunkExists(uploadID string, chunkIndex int) (bool, error)
	MergeChunks(uploadID string, chunkLength int, finalFileName string) error
	CleanupChunks(uploadID string) error
	DeleteFinalFile(fileName string) error
}
