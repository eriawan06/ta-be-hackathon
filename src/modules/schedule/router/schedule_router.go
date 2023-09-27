package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/schedule"
	"be-sagara-hackathon/src/utils/constants"
	"github.com/gin-gonic/gin"
)

func ScheduleRouter(group *gin.RouterGroup) {
	group.POST("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		schedule.GetController().CreateSchedule,
	)
	group.POST("/teams",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		schedule.GetController().CreateScheduleTeam,
	)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		schedule.GetController().UpdateSchedule,
	)
	group.DELETE("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		schedule.GetController().DeleteSchedule,
	)
	group.DELETE("/:id/teams/:team_id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		schedule.GetController().DeleteScheduleTeam,
	)
	group.GET("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin, constants.UserMentor),
		schedule.GetController().GetListSchedule,
	)
	group.GET("/:id", schedule.GetController().GetDetailSchedule)
}
