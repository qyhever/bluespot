package model

type AttachUploadResponse struct {
	FileName   string `json:"fileName"`
	OriginName string `json:"originName"`
	URL        string `json:"url"`
}
