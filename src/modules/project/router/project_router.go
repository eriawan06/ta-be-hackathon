package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/project"
	"be-sagara-hackathon/src/utils/constants"

	"github.com/gin-gonic/gin"
)

func ProjectRouter(group *gin.RouterGroup) {
	group.POST("/",
		middlewares.RolePermission(constants.UserParticipant),
		project.GetProjectController().Create,
	)
	group.POST("/:id/assessments",
		middlewares.RolePermission(constants.UserJudge),
		project.GetProjectAssessmentController().Create,
	)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserParticipant),
		project.GetProjectController().Update,
	)
	group.PUT("/:id/status/:status",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		project.GetProjectController().UpdateStatus,
	)
	group.GET("/:id", project.GetProjectController().GetDetail)
	group.GET("/:id/assessments",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		project.GetProjectAssessmentController().GetByProjectID,
	)
	group.GET("/:id/assessments/judge",
		middlewares.RolePermission(constants.UserJudge),
		project.GetProjectAssessmentController().GetByJudgeAndProjectID,
	)
	group.GET("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin, constants.UserJudge),
		project.GetProjectController().GetAll,
	)
}
