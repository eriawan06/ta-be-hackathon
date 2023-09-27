package controller

import (
	"be-sagara-hackathon/src/modules/auth/model"
	"be-sagara-hackathon/src/modules/auth/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	"be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminAuthController interface {
	Login(ctx *gin.Context)
}

type AdminAuthControllerImpl struct {
	Service service.AdminAuthService
}

func NewAdminAuthController(service service.AdminAuthService) AdminAuthController {
	return &AdminAuthControllerImpl{Service: service}
}

// Login Login Admin godoc
// @Tags Authentication
// @Summary Login Admin
// @Description Login For Superadmin, Admin, HR
// @Accept  json
// @Produce  json
// @Param body body dto.LoginRequest true "Body Request"
// @Success 200 {object} src.AuthSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /auth/login/admin [post]
func (controller *AdminAuthControllerImpl) Login(ctx *gin.Context) {
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
		if err == errors.ErrWrongLoginCredential {
			common.SendError(ctx, http.StatusUnauthorized, "Unauthorized", []string{err.Error()})
			return
		}

		if err == errors.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Login Success", &response)
}
