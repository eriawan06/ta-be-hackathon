package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/event"
	"be-sagara-hackathon/src/utils/constants"

	"github.com/gin-gonic/gin"
)

func EventRouter(group *gin.RouterGroup) {
	group.POST("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		event.GetController().Create,
	)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		event.GetController().Update,
	)
	group.DELETE("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		event.GetController().Delete,
	)
	group.GET("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		event.GetController().GetList,
	)
	group.GET("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		event.GetController().GetDetail,
	)
	group.GET("/:id/rules",
		middlewares.RolePermission(constants.UserParticipant),
		event.GetEventRuleController().GetActiveByEventID,
	)
	group.GET("/latest", event.GetController().GetLatest)
	group.GET("/:id/schedules",
		middlewares.RolePermission(constants.UserParticipant),
		event.GetController().GetSchedules,
	)

	/// Event Mentor Routes ///
	em := group.Group("/mentors")
	{
		em.POST("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventMentorController().Create,
		)
		em.DELETE("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventMentorController().Delete,
		)
		em.GET("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventMentorController().GetAll,
		)
		em.GET("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventMentorController().GetDetail,
		)
	}

	/// Event Judge Routes ///
	ej := group.Group("/judges")
	{
		ej.POST("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventJudgeController().Create,
		)
		ej.DELETE("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventJudgeController().Delete,
		)
		ej.GET("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventJudgeController().GetAll,
		)
		ej.GET("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventJudgeController().GetDetail,
		)
	}

	/// Event Company Routes ///
	ec := group.Group("/companies")
	{
		ec.POST("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventCompanyController().Create,
		)
		ec.PUT("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventCompanyController().Update,
		)
		ec.DELETE("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventCompanyController().Delete,
		)
		ec.GET("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventCompanyController().GetList,
		)
		ec.GET("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventCompanyController().GetDetail,
		)
	}

	/// Event Timeline Routes ///
	etl := group.Group("/timelines")
	{
		etl.POST("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventTimelineController().Create,
		)
		etl.PUT("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventTimelineController().Update,
		)
		etl.DELETE("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventTimelineController().Delete,
		)
		etl.GET("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventTimelineController().GetList,
		)
		etl.GET("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventTimelineController().GetDetail,
		)
	}

	/// Event Rules Routes ///
	evr := group.Group("/rules")
	{
		evr.POST("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventRuleController().Create,
		)
		evr.PUT("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventRuleController().Update,
		)
		evr.DELETE("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventRuleController().Delete,
		)
		evr.GET("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventRuleController().GetList,
		)
		evr.GET("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventRuleController().GetDetail,
		)
	}

	/// Event FAQ Routes ///
	efaq := group.Group("/faqs")
	{
		efaq.POST("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventFaqController().Create,
		)
		efaq.PUT("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventFaqController().Update,
		)
		efaq.DELETE("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventFaqController().Delete,
		)
		efaq.GET("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventFaqController().GetList,
		)
		efaq.GET("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventFaqController().GetDetail,
		)
	}

	/// Event Assessment Criteria Routes ///
	eac := group.Group("/assessments")
	{
		eac.POST("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventAssessmentCriteriaController().Create,
		)
		eac.PUT("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventAssessmentCriteriaController().Update,
		)
		eac.DELETE("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventAssessmentCriteriaController().Delete,
		)
		eac.GET("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin, constants.UserJudge),
			event.GetEventAssessmentCriteriaController().GetList,
		)
		eac.GET("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			event.GetEventAssessmentCriteriaController().GetDetail,
		)
	}
}
