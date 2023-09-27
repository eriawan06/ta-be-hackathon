package router

import (
	"be-sagara-hackathon/src/modules/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(group *gin.RouterGroup) {
	authController := auth.GetController()
	group.POST("/register", authController.Register)
	group.POST("/login", authController.Login)
	group.GET("/google/callback", authController.GoogleOauth)
	group.POST("/register/google", authController.RegisterByGoogle)
	group.POST("/login/google", authController.LoginByGoogle)
	group.POST("/verify-email", authController.VerifyEmail)
	group.POST("/get-verification-code", authController.SendVerificationCode)
	group.POST("/validate-verification-code", authController.ValidateVerificationCode)
	group.POST("/forgot-password", authController.ForgotPassword)
	//group.POST("/reset-password", authController.ResetPassword) //DEPRECATED

	//admin & internal
	adminAuthController := auth.GetAdminController()
	group.POST("/login/admin", adminAuthController.Login)
}
