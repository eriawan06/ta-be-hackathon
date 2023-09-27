package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/payment"
	"be-sagara-hackathon/src/utils/constants"
	"github.com/gin-gonic/gin"
)

func PaymentMethodRouter(group *gin.RouterGroup) {
	group.POST("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		payment.GetPaymentMethodController().Create,
	)
	group.PUT("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		payment.GetPaymentMethodController().Update,
	)
	group.DELETE("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		payment.GetPaymentMethodController().Delete,
	)
	group.GET("/", payment.GetPaymentMethodController().GetList)
	group.GET("/:id", payment.GetPaymentMethodController().GetOne)
}
