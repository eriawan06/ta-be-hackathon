package router

import (
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/payment"
	"be-sagara-hackathon/src/utils/constants"
	"github.com/gin-gonic/gin"
)

func InvoiceRouter(group *gin.RouterGroup) {
	group.GET("/",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		payment.GetInvoiceController().GetList,
	)
	group.GET("/:id",
		middlewares.RolePermission(constants.UserSuperadmin, constants.UserAdmin),
		payment.GetInvoiceController().GetDetail,
	)
	group.GET("/:id/payments", payment.GetPaymentController().GetManyByInvoiceID)
}
