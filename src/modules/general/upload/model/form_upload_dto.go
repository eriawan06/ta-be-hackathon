package model

import (
	"be-sagara-hackathon/src/utils/upload"
	"mime/multipart"
)

type FormUploadRequest struct {
	File      multipart.File
	FileInfo  upload.FileInfo
	Path      string `form:"path" binding:"required"`
	Overwrite bool   `form:"overwrite"`
	PrevFile  string `form:"previous_file"`
}

type UploadResponse struct {
	FileUrl  string `json:"file_url"`
	FilePath string `json:"file_path"`
}
