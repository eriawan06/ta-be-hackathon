package router

import (
	"be-sagara-hackathon/src/modules/general/upload"
	"github.com/gin-gonic/gin"
)

func UploadRouter(group *gin.RouterGroup) {
	group.POST("/", upload.GetUploadController().Upload)
}
