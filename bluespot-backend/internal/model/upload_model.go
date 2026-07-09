package model

type UploadVerifyRequest struct {
	FileMD5  string `json:"fileMd5"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
}

type UploadVerifyResponse struct {
	IsExists       bool   `json:"isExists"`
	URL            string `json:"url"`
	UploadID       string `json:"uploadId"`
	UploadedChunks []int  `json:"uploadedChunks"`
}

type UploadMergeRequest struct {
	UploadID    string `json:"uploadId"`
	FileMD5     string `json:"fileMd5"`
	ChunkLength int    `json:"chunkLength"`
}

type UploadMergeResponse struct {
	URL string `json:"url"`
	Msg string `json:"msg"`
}
