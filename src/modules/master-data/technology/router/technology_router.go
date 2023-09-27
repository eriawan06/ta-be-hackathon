package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/master-data/technology"
	"be-sagara-hackathon/src/utils/constants"
	"github.com/gin-gonic/gin"
)

func TechnologyRouter(group *gin.RouterGroup) {
	group.GET("/", technology.GetController().GetListTechnology)
	group.GET("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		technology.GetController().GetDetailTechnology,
	)
	group.POST("/", technology.GetController().CreateTechnology)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		technology.GetController().UpdateTechnology,
	)
}
