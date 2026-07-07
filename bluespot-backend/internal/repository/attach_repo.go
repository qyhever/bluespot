package repository

import "mime/multipart"

type AttachRepository interface {
	Save(file *multipart.FileHeader, fileName string) error
}
