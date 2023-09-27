package controller

import (
	"be-sagara-hackathon/src/modules/general/upload/model"
	"be-sagara-hackathon/src/modules/general/upload/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/upload"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
)

type UploadController interface {
	Upload(ctx *gin.Context)
}

type UploadControllerImplS3 struct {
	Service service.UploadService
}

func NewUploadControllerS3(uploadService service.UploadService) UploadController {
	return &UploadControllerImplS3{
		Service: uploadService,
	}
}

// Upload UploadFile godoc
// @Tags General
// @Summary Upload File
// @Description supported format : ".png", ".jpg", ".jpeg", ".pdf"
// @Description list path : "avatars/participants", "avatars/mentors", "avatars/judges", "avatars/admin", "cv", "payment-evidence"
// @Description Field overwrite is used for updating file. Set overwrite to true and fill previous_file when you want to update file
// @Accept multipart/form-data
// @Produce  json
// @Security ApiKeyAuth
// @Param body body src.UploadRequest true "Body Request"
// @Success 200 {object} src.UploadFileSuccess
// @Success 400 {object} src.BaseFailure
// @Router /upload [post]
func (controller UploadControllerImplS3) Upload(ctx *gin.Context) {
	_, err := utils.GetUserCredentialFromToken(ctx)
	if err != nil {
		common.SendError(ctx, http.StatusUnauthorized, "Invalid Token", []string{err.Error()})
		return
	}

	// limit upload file size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 2*upload.MB) // 2 Mb

	var request model.FormUploadRequest
	errorBinding := ctx.ShouldBind(&request)
	if errorBinding != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{errorBinding.Error()})
		return
	}

	var multipartFileHeader *multipart.FileHeader
	request.File, multipartFileHeader, err = ctx.Request.FormFile("file")
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	// Create a buffer to store the header of the file in
	// And copy the headers into the FileHeader buffer
	fileHeader := make([]byte, 512)
	if _, err := request.File.Read(fileHeader); err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	// set position back to start.
	if _, err := request.File.Seek(0, 0); err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	defer request.File.Close()

	mType := mimetype.Detect(fileHeader)
	request.FileInfo = upload.FileInfo{
		FileName: multipartFileHeader.Filename,
		FileSize: request.File.(upload.Sizer).Size(),
		FileMime: mType.String(),
		FileExt:  mType.Extension(),
	}

	data, err := controller.Service.UploadToS3(request)
	if err != nil {
		if err == e.ErrUnsupportedFileFormat ||
			err == e.ErrWrongFileUploadPath {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Upload File Success", data)
}
