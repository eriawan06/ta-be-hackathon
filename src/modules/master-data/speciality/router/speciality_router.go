package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/master-data/speciality"
	"be-sagara-hackathon/src/utils/constants"
	"github.com/gin-gonic/gin"
)

func SpecialityRouter(group *gin.RouterGroup) {
	group.POST("/", speciality.GetController().CreateSpeciality)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		speciality.GetController().UpdateSpeciality,
	)
	group.GET("/", speciality.GetController().GetListSpeciality)
	group.GET("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		speciality.GetController().GetDetailSpeciality,
	)
}
