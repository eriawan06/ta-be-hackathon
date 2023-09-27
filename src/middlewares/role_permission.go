package middlewares

import (
	"be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"be-sagara-hackathon/src/utils/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RolePermission(roles ...string) gin.HandlerFunc {
	return func(context *gin.Context) {

		userExtract, _ := context.Get("user")
		user := userExtract.(model.User)

		// check permission
		if len(roles) > 0 && !helper.StringInSlice(user.UserRole.Name, roles) {
			common.SendError(context, http.StatusForbidden, "Forbidden", []string{"Forbidden"})
			context.Abort()
			return
		}

		context.Next()
	}
}
