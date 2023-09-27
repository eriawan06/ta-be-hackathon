package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/master-data/occupation"
	"be-sagara-hackathon/src/utils/constants"
	"github.com/gin-gonic/gin"
)

func OccupationRouter(group *gin.RouterGroup) {
	group.POST("/", occupation.GetController().CreateOccupation)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		occupation.GetController().UpdateOccupation,
	)
	group.GET("/", occupation.GetController().GetListOccupation)
	group.GET("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		occupation.GetController().GetDetailOccupation,
	)
}
