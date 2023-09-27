package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/payment"
	"be-sagara-hackathon/src/modules/user"
	"be-sagara-hackathon/src/utils/constants"

	"github.com/gin-gonic/gin"
)

func UserRouter(group *gin.RouterGroup) {
	/// User Routes ///
	group.POST("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		user.GetUserController().CreateUser,
	)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		user.GetUserController().UpdateUser,
	)
	group.DELETE("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		user.GetUserController().DeleteUser,
	)
	group.GET("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		user.GetUserController().GetList,
	)
	group.GET("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		user.GetUserController().GetDetail,
	)
	group.GET("/profile", user.GetUserController().GetUserProfile)
	group.PUT("/change-password", user.GetUserController().ChangePassword)

	/// Participant Routes ///
	participant := group.Group("/participants")
	{
		participant.GET("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetParticipantController().GetList,
		)
		participant.GET("/detail/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetParticipantController().GetDetail,
		)
		participant.GET("/profile",
			middlewares.RolePermission(constants.UserParticipant),
			user.GetParticipantController().GetProfile,
		)
		participant.GET("/:participant_id/events/:event_id/invoice",
			middlewares.RolePermission(constants.UserParticipant),
			payment.GetInvoiceController().GetParticipantInvoice,
		)
		participant.PUT("/",
			middlewares.RolePermission(constants.UserParticipant),
			user.GetParticipantController().UpdateFull,
		)
		participant.PUT("/profile",
			middlewares.RolePermission(constants.UserParticipant),
			user.GetParticipantController().UpdateProfileAndLocation,
		)
		participant.PUT("/education",
			middlewares.RolePermission(constants.UserParticipant),
			user.GetParticipantController().UpdateEducation,
		)
		participant.PUT("/preference",
			middlewares.RolePermission(constants.UserParticipant),
			user.GetParticipantController().UpdatePreference,
		)
		participant.PUT("/account",
			middlewares.RolePermission(constants.UserParticipant),
			user.GetParticipantController().UpdateAccount,
		)
		participant.POST("/complete-registration",
			middlewares.RolePermission(constants.UserParticipant),
			user.GetParticipantController().CompleteRegistration,
		)
		participant.POST("/search",
			middlewares.RolePermission(constants.UserParticipant),
			user.GetParticipantController().GetListParticipantSearch,
		)
	}

	/// Mentor Routes ///
	mentor := group.Group("/mentors")
	{
		mentor.POST("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetMentorController().Create,
		)
		mentor.PUT("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetMentorController().Update,
		)
		mentor.DELETE("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetMentorController().Delete,
		)
		mentor.GET("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetMentorController().GetList,
		)
		mentor.GET("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetMentorController().GetDetail,
		)
	}

	/// Judge Routes ///
	judge := group.Group("/judges")
	{
		judge.POST("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetJudgeController().Create,
		)
		judge.PUT("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetJudgeController().Update,
		)
		judge.DELETE("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetJudgeController().Delete,
		)
		judge.GET("/",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetJudgeController().GetList,
		)
		judge.GET("/:id",
			middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
			user.GetJudgeController().GetDetail,
		)
	}
}
