package service

import (
	"be-sagara-hackathon/src/modules/general/upload/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
	"be-sagara-hackathon/src/utils/upload"
	"fmt"
	"os"
)

var (
	listAllowedFormatFile = []string{".png", ".jpg", ".jpeg", ".pdf", ".zip"}
	listPath              = []string{
		"avatars/participants", "avatars/mentors", "avatars/judges", "avatars/admin",
		"cv", "payment-evidence", "logo/companies", "logo/teams",
		"project/thumbnail", "project/images",
	}
)

type UploadService interface {
	UploadToS3(upload model.FormUploadRequest) (model.UploadResponse, error)
}

type UploadServiceImpl struct{}

func NewUploadService() UploadService {
	return &UploadServiceImpl{}
}

func (service UploadServiceImpl) UploadToS3(formUpload model.FormUploadRequest) (model.UploadResponse, error) {
	//check file format or extension
	if !helper.StringInSlice(formUpload.FileInfo.FileExt, listAllowedFormatFile) {
		return model.UploadResponse{}, e.ErrUnsupportedFileFormat
	}

	//check path
	if !helper.StringInSlice(formUpload.Path, listPath) {
		return model.UploadResponse{}, e.ErrWrongFileUploadPath
	}

	formUpload.FileInfo.FileName = formUpload.PrevFile
	if !formUpload.Overwrite {
		formUpload.FileInfo.FileName = fmt.Sprintf(
			"%s/%s%s",
			formUpload.Path,
			utils.GenerateRandomAlphaNumberic(10),
			formUpload.FileInfo.FileExt)
	}

	//// upload to s3
	err := upload.PushS3(upload.S3Info{
		//Endpoint: helper.ReferString(os.Getenv("LINODE_ENDPOINT")),
		Key:      os.Getenv("AWS_S3_ACCESS_KEY"),
		Secret:   os.Getenv("AWS_S3_SECRET_KEY"),
		Region:   os.Getenv("AWS_S3_REGION"),
		Bucket:   os.Getenv("AWS_S3_BUCKET"),
		File:     formUpload.File,
		Filename: formUpload.FileInfo.FileName,
		Filemime: formUpload.FileInfo.FileMime,
		Filesize: formUpload.FileInfo.FileSize,
	})
	if err != nil {
		return model.UploadResponse{}, err
	}

	return model.UploadResponse{
		FileUrl:  fmt.Sprintf("%s/%s", os.Getenv("LINODE_BASE_FILE_URL"), formUpload.FileInfo.FileName),
		FilePath: formUpload.FileInfo.FileName,
	}, nil
}
