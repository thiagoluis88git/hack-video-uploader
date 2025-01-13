package entity

import "mime/multipart"

type UoloaderDocumentEntity struct {
	Data        multipart.File
	Name        string
	ContentType string
	Size        int64
}
