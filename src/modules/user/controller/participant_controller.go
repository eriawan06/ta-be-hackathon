package controller

import (
	"be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/modules/user/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ParticipantController interface {
	GetProfile(ctx *gin.Context)
	UpdateFull(ctx *gin.Context)
	UpdateProfileAndLocation(ctx *gin.Context)
	UpdateEducation(ctx *gin.Context)
	UpdatePreference(ctx *gin.Context)
	UpdateAccount(ctx *gin.Context)
	CompleteRegistration(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	GetListParticipantSearch(ctx *gin.Context)
}

type ParticipantControllerImpl struct {
	Service service.ParticipantService
}

func NewParticipantController(service service.ParticipantService) ParticipantController {
	return &ParticipantControllerImpl{Service: service}
}

func (controller *ParticipantControllerImpl) GetProfile(ctx *gin.Context) {
	profile, err := controller.Service.GetProfile(ctx)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Participant Profile Success", profile)
}

func (controller *ParticipantControllerImpl) UpdateFull(ctx *gin.Context) {
	var request model.UpdateParticipantFull
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		if err.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		// When Binding Error
		common.SendError(ctx, http.StatusBadRequest, "Not valid request", utils.SplitError(err))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	data, err := controller.Service.UpdateFull(ctx, request)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Participant Success", data)
}

func (controller *ParticipantControllerImpl) UpdateProfileAndLocation(ctx *gin.Context) {
	var request model.UpdateParticipantProfileRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		if err.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		// When Binding Error
		common.SendError(ctx, http.StatusBadRequest, "Not valid request", utils.SplitError(err))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	data, err := controller.Service.UpdateProfileAndLocation(ctx, request)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Participant Profile & Location Success", data)
}

func (controller *ParticipantControllerImpl) UpdateEducation(ctx *gin.Context) {
	var request model.UpdateParticipantEducationRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		if err.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		// When Binding Error
		common.SendError(ctx, http.StatusBadRequest, "Not valid request", utils.SplitError(err))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	data, err := controller.Service.UpdateEducation(ctx, request)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Participant Education Success", data)
}

func (controller *ParticipantControllerImpl) UpdatePreference(ctx *gin.Context) {
	var request model.UpdateParticipantPreferenceRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		if err.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		// When Binding Error
		common.SendError(ctx, http.StatusBadRequest, "Not valid request", utils.SplitError(err))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	data, err := controller.Service.UpdatePreference(ctx, request)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Participant Preference Success", data)
}

func (controller *ParticipantControllerImpl) UpdateAccount(ctx *gin.Context) {
	var request model.UpdateParticipantAccountRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		if err.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		// When Binding Error
		common.SendError(ctx, http.StatusBadRequest, "Not valid request", utils.SplitError(err))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	data, err := controller.Service.UpdateAccount(ctx, request)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		if err == e.ErrEmailAlreadyExists || err == e.ErrPhoneNumberAlreadyExists || err == e.ErrUsernameAlreadyExists {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Participant Account Success", data)
}

func (controller *ParticipantControllerImpl) CompleteRegistration(ctx *gin.Context) {
	err := controller.Service.CompleteRegistration(ctx)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		if err == e.ErrRegistrationNotCompleted {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Complete Registration Success", nil)
}

func (controller *ParticipantControllerImpl) GetList(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	var filter model.FilterUser
	filter.Status = ctx.Query("status")
	filter.Search = ctx.Query("q")

	data, err := controller.Service.GetList(filter, pg)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Participant Success", data)
}

func (controller *ParticipantControllerImpl) GetDetail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	profile, err := controller.Service.GetDetail(uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Participant Success", profile)
}

func (controller *ParticipantControllerImpl) GetListParticipantSearch(ctx *gin.Context) {
	var filter model.FilterParticipantSearch
	err := ctx.ShouldBindJSON(&filter)
	if err != nil {
		if err.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		// When Binding Error
		common.SendError(ctx, http.StatusBadRequest, "Invalid filter", utils.SplitError(err))
		return
	}

	// Validate filter body
	if errs := utils.NewCustomValidator().ValidateStruct(filter); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid filter", errs)
		return
	}

	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetListParticipantSearch(filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Participant Search Success", data)
}
