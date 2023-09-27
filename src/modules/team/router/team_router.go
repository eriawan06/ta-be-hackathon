package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/team"
	"be-sagara-hackathon/src/utils/constants"
	"github.com/gin-gonic/gin"
)

func TeamRouter(group *gin.RouterGroup) {
	group.POST("/",
		middlewares.RolePermission(constants.UserParticipant),
		team.GetTeamController().Create,
	)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserParticipant),
		team.GetTeamController().Update,
	)
	group.PUT("/:id/status",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		team.GetTeamController().UpdateStatus,
	)
	group.DELETE("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		team.GetTeamController().Delete,
	)
	group.GET("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin, constants.UserMentor),
		team.GetTeamController().GetAll,
	)
	group.GET("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		team.GetTeamController().GetDetail,
	)
	group.GET("/:id/detail",
		middlewares.RolePermission(constants.UserParticipant),
		team.GetTeamController().GetDetail2,
	)
	group.GET("/:id/members",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin, constants.UserParticipant),
		team.GetTeamController().GetMembers,
	)
	group.GET("/:id/invitations",
		middlewares.RolePermission(constants.UserParticipant),
		team.GetTeamController().GetInvitations,
	)
	group.GET("/:id/requests",
		middlewares.RolePermission(constants.UserParticipant),
		team.GetTeamController().GetRequests,
	)
	group.GET("/my-team",
		middlewares.RolePermission(constants.UserParticipant),
		team.GetTeamController().GetMyTeam,
	)
	group.GET("/events/:event_id",
		middlewares.RolePermission(constants.UserParticipant),
		team.GetTeamController().GetListByEventID,
	)

	member := group.Group("/members")
	{
		member.DELETE("/:id",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamMemberController().Delete,
		)
	}

	invitation := group.Group("/invitations")
	{
		invitation.POST("/",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamInvitationController().Create,
		)
		invitation.PUT("/:id",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamInvitationController().Update,
		)
		invitation.PUT("/:id/status",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamInvitationController().UpdateStatus,
		)
		invitation.DELETE("/:id",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamInvitationController().Delete,
		)
		invitation.GET("/",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamInvitationController().GetList,
		)
		invitation.GET("/:id",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamInvitationController().GetDetail,
		)
		invitation.GET("/:id/detail",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamInvitationController().GetDetail2,
		)
	}

	request := group.Group("/requests")
	{
		request.POST("/",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamRequestController().Create,
		)
		request.PUT("/:id",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamRequestController().Update,
		)
		request.DELETE("/:id",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamRequestController().Delete,
		)
		request.GET("/:id",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamRequestController().GetDetail,
		)
		request.PUT("/:id/status",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamRequestController().UpdateStatus,
		)
		request.GET("/:id/detail",
			middlewares.RolePermission(constants.UserParticipant),
			team.GetTeamRequestController().GetDetailFull,
		)
	}
}
