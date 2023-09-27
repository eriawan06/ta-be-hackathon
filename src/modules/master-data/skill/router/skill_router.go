package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/master-data/skill"
	"be-sagara-hackathon/src/utils/constants"
	"github.com/gin-gonic/gin"
)

func SkillRouter(group *gin.RouterGroup) {
	group.GET("/", skill.GetController().GetListSkill)
	group.GET("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		skill.GetController().GetDetailSkill,
	)
	group.POST("/", skill.GetController().CreateSkill)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		skill.GetController().UpdateSkill,
	)
}
