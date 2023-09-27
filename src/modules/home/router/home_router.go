package router

import (
	"be-sagara-hackathon/src/modules/home"
	"github.com/gin-gonic/gin"
)

func HomeRouter(group *gin.RouterGroup) {
	group.GET("/", home.GetController().GetData)
}
