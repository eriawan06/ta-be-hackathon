package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/payment"
	"be-sagara-hackathon/src/utils/constants"
	"github.com/gin-gonic/gin"
)

func PaymentRouter(group *gin.RouterGroup) {
	group.POST("/",
		middlewares.RolePermission(constants.UserParticipant),
		payment.GetPaymentController().Create,
	)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		payment.GetPaymentController().Update,
	)
	group.GET("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		payment.GetPaymentController().GetList,
	)
	group.GET("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		payment.GetPaymentController().GetDetail,
	)
}
