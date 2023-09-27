package upload

import (
	"be-sagara-hackathon/src/modules/general/upload/controller"
	"be-sagara-hackathon/src/modules/general/upload/service"
)

var (
	uploadController controller.UploadController
	uploadService    service.UploadService
)

type Module interface {
	InitModule()
}

type ModuleImpl struct{}

func New() Module {
	return &ModuleImpl{}
}

func (module ModuleImpl) InitModule() {
	uploadService = service.NewUploadService()
	uploadController = controller.NewUploadControllerS3(uploadService)
}

func GetUploadController() controller.UploadController {
	return uploadController
}
