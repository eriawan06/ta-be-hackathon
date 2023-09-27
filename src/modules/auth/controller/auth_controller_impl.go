package controller

import (
	_ "be-sagara-hackathon/src"
	"be-sagara-hackathon/src/modules/auth/model"
	"be-sagara-hackathon/src/modules/auth/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/oauth"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	RegisterByGoogle(ctx *gin.Context)
	LoginByGoogle(ctx *gin.Context)
	GoogleOauth(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	SendVerificationCode(ctx *gin.Context)
	ValidateVerificationCode(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
}

// AuthControllerImpl Binding Services to Controller
type AuthControllerImpl struct {
	Service service.AuthService
}

// NewAuthController Create Function to Create New Instance
func NewAuthController(authService service.AuthService) AuthController {
	return &AuthControllerImpl{
		Service: authService,
	}
}

// Register Register By Google godoc
// @Tags Authentication
// @Summary Register Regular
// @Description Register User Regular
// @Accept  json
// @Produce  json
// @Param body body dto.RegisterRequest true "Body Request"
// @Success 200 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /auth/register [post]
func (controller *AuthControllerImpl) Register(ctx *gin.Context) {
	// Create New Request Model
	var request model.RegisterRequest
	if errorBinding := ctx.ShouldBindJSON(&request); errorBinding != nil {
		if errorBinding.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(errorBinding))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	err := controller.Service.Register(request)
	if err != nil {
		if err == e.ErrEmailAlreadyExists ||
			err == e.ErrPhoneNumberAlreadyExists ||
			err == e.ErrConfirmPasswordNotSame {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Register Success", nil)
}

// Login Login By Google godoc
// @Tags Authentication
// @Summary Login Regular
// @Description Login User Regular
// @Accept  json
// @Produce  json
// @Param body body model.LoginRequest true "Body Request"
// @Success 200 {object} src.AuthSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /auth/login [post]
func (controller *AuthControllerImpl) Login(ctx *gin.Context) {
	var request model.LoginRequest
	if errorBinding := ctx.ShouldBindJSON(&request); errorBinding != nil {
		if errorBinding.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(errorBinding))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	response, err := controller.Service.Login(request)
	if err != nil {
		if err == e.ErrWrongLoginCredential || err == e.ErrUserIsNotActivated {
			common.SendError(ctx, http.StatusUnauthorized, "Unauthorized", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Login Success", &response)
}

// RegisterByGoogle Register By Google godoc
// @Tags Authentication
// @Summary Register By Google
// @Description Register User By Google Account
// @Accept  json
// @Produce  json
// @Param body body model.RegisterByGoogleRequest true "Body Request"
// @Success 200 {object} src.AuthSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /auth/register/google [post]
func (controller *AuthControllerImpl) RegisterByGoogle(ctx *gin.Context) {
	var request model.RegisterByGoogleRequest
	if errorBinding := ctx.ShouldBindJSON(&request); errorBinding != nil {
		if errorBinding.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(errorBinding))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	response, err := controller.Service.RegisterByGoogle(request)
	if err != nil {
		// When Id Token Expired or User Not Found
		if err.Error() == "idtoken: token expired" {
			common.SendError(ctx, http.StatusUnauthorized, "Unauthorized", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Register Success", &response)
}

// LoginByGoogle Login By Google godoc
// @Tags Authentication
// @Summary Login By Google
// @Description Login User By Google Account
// @Accept  json
// @Produce  json
// @Param body body model.LoginByGoogleRequest true "Body Request"
// @Success 200 {object} src.AuthSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /auth/login/google [post]
func (controller *AuthControllerImpl) LoginByGoogle(ctx *gin.Context) {
	var request model.LoginByGoogleRequest
	if errorBinding := ctx.ShouldBindJSON(&request); errorBinding != nil {
		if errorBinding.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(errorBinding))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	response, err := controller.Service.LoginByGoogle(request)
	if err != nil {
		// When Id Token Expired
		if (err.Error() == "idtoken: token expired") || (err.Error() == "user not found") {
			common.SendError(ctx, http.StatusUnauthorized, "Unauthorized", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Login Success", &response)
}

func (controller *AuthControllerImpl) GoogleOauth(ctx *gin.Context) {
	redirectUrlPath := os.Getenv("OAUTH_GOOGLE_REDIRECT_URL")
	redirectUrl := fmt.Sprintf("%s%s", os.Getenv("BASE_FE_URL"), redirectUrlPath)

	code := ctx.Query("code")
	if code == "" {
		ctx.SetCookie("is_authenticated", "false", 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), false, false)
		ctx.SetCookie("error", "Authorization code not provided!", 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), true, true)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		//common.SendError(ctx, http.StatusUnauthorized, "Unauthorized", []string{"Authorization code not provided!"})
		return
	}

	// Use the code to get the id and access tokens
	tokenRes, err := oauth.GetGoogleOauthToken(code)
	if err != nil {
		ctx.SetCookie("is_authenticated", "false", 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), false, false)
		ctx.SetCookie("error", err.Error(), 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), true, true)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		//common.SendError(ctx, http.StatusBadGateway, "Bad Gateway", []string{err.Error()})
		return
	}

	googleUser, err := oauth.GetGoogleUser(tokenRes.AccessToken, tokenRes.IdToken)
	if err != nil {
		ctx.SetCookie("is_authenticated", "false", 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), false, false)
		ctx.SetCookie("error", err.Error(), 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), true, true)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		//common.SendError(ctx, http.StatusBadGateway, "Bad Gateway", []string{err.Error()})
		return
	}

	response, err := controller.Service.GoogleOauth(*googleUser)
	if err != nil {
		ctx.SetCookie("is_authenticated", "false", 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), false, false)
		ctx.SetCookie("error", err.Error(), 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), true, true)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}

	userData, err := json.Marshal(response.User)
	if err != nil {
		ctx.SetCookie("is_authenticated", "false", 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), false, false)
		ctx.SetCookie("error", err.Error(), 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), true, true)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}

	ctx.SetCookie("access_token", response.Token, 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), false, true)
	ctx.SetCookie("user", string(userData), 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), false, true)
	ctx.SetCookie("is_authenticated", "true", 1000*60, redirectUrlPath, os.Getenv("BASE_FE_URL"), false, false)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

// VerifyEmail Email Verification godoc
// @Tags Authentication
// @Summary Email Verification
// @Description Email Verification
// @Accept  json
// @Produce  json
// @Param body body model.VerifyEmailRequest true "Body Request"
// @Success 200 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /auth/verify-email [post]
func (controller *AuthControllerImpl) VerifyEmail(ctx *gin.Context) {
	var request model.VerifyEmailRequest
	if errorBinding := ctx.ShouldBindJSON(&request); errorBinding != nil {
		if errorBinding.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(errorBinding))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	err := controller.Service.VerifyEmail(request)
	if err != nil {
		if err == e.ErrEmailNotRegistered ||
			err == e.ErrEmailAlreadyVerified ||
			err == e.ErrInvalidVerificationCode {
			common.SendError(ctx, http.StatusBadRequest, "Invalid request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Verify Email Success", nil)
}

// SendVerificationCode Get Verification Code godoc
// @Tags Authentication
// @Summary Get Verification Code
// @Description Get Verification Code
// @Accept  json
// @Produce  json
// @Param body body model.SendVerificationCodeRequest true "Body Request"
// @Success 200 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /auth/get-verification-code [post]
func (controller *AuthControllerImpl) SendVerificationCode(ctx *gin.Context) {
	var request model.SendVerificationCodeRequest
	if errorBinding := ctx.ShouldBindJSON(&request); errorBinding != nil {
		if errorBinding.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(errorBinding))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	err := controller.Service.SendVerificationCode(request)
	if err != nil {
		if err == e.ErrEmailNotRegistered {
			common.SendError(ctx, http.StatusBadRequest, "Invalid request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Email Verification Code Success", nil)
}

// ValidateVerificationCode Validate Verification Code godoc
// @Tags Authentication
// @Summary Validate Verification Code
// @Description Validate Verification Code
// @Accept  json
// @Produce  json
// @Param body body model.ValidateVerificationCodeRequest true "Body Request"
// @Success 200 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /auth/validate-verification-code [post]
func (controller *AuthControllerImpl) ValidateVerificationCode(ctx *gin.Context) {
	var request model.ValidateVerificationCodeRequest
	if errorBinding := ctx.ShouldBindJSON(&request); errorBinding != nil {
		if errorBinding.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(errorBinding))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	err := controller.Service.ValidateVerificationCode(request)
	if err != nil {
		if err == e.ErrInvalidVerificationCode {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Validate Verification Code Success", nil)
}

// ForgotPassword Forgot Password godoc
// @Tags Authentication
// @Summary Forgot Password
// @Description Forgot Password
// @Accept  json
// @Produce  json
// @Param body body model.ForgotPasswordRequest true "Body Request"
// @Success 200 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /auth/forgot-password [post]
func (controller *AuthControllerImpl) ForgotPassword(ctx *gin.Context) {
	var request model.ForgotPasswordRequest
	if errorBinding := ctx.ShouldBindJSON(&request); errorBinding != nil {
		if errorBinding.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(errorBinding))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	err := controller.Service.ForgotPassword(request)
	if err != nil {
		if err == e.ErrInvalidVerificationCode {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Forgot Password Success", nil)
}
