package controller

import (
	"be-sagara-hackathon/src/modules/team/model"
	"be-sagara-hackathon/src/modules/team/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TeamInvitationController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	UpdateStatus(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	GetDetail2(ctx *gin.Context)
}

type TeamInvitationControllerImpl struct {
	Service service.TeamInvitationService
}

func NewTeamInvitationController(service service.TeamInvitationService) TeamInvitationController {
	return &TeamInvitationControllerImpl{Service: service}
}

func (controller *TeamInvitationControllerImpl) Create(ctx *gin.Context) {
	var request model.CreateInvitationRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		if err.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(err))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	err = controller.Service.Create(ctx, request)
	if err != nil {
		if err == e.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrEventNotRunning || err == e.ErrHasTeam || err == e.ErrParticipantHasBeenInvited ||
			err == e.ErrRegistrationNotCompleted || err == e.ErrPaymentNotPaid || err == e.ErrTeamIsFull ||
			err == e.ErrParticipantRequestedToJoinTeam {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Team Invitation Success", nil)
}

func (controller *TeamInvitationControllerImpl) Update(ctx *gin.Context) {
	var request model.UpdateInvitationRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		if err.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(err))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	err = controller.Service.Update(ctx, uint(id), request)
	if err != nil {
		if err == e.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrInvitationHasBeenProceed {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Team Invitation Success", nil)
}

func (controller *TeamInvitationControllerImpl) UpdateStatus(ctx *gin.Context) {
	var request model.UpdateStatusInvitationRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		if err.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(err))
		return
	}

	// Validate request body
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	code := ctx.Param("id")
	if code == "" {
		common.SendError(ctx, http.StatusBadRequest, "Invalid code", []string{err.Error()})
		return
	}

	err = controller.Service.UpdateStatus(ctx, code, request)
	if err != nil {
		if err == e.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrInvitationHasBeenProceed || err == e.ErrHasTeam || err == e.ErrTeamIsFull {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Team Invitation Status Success", nil)
}

func (controller *TeamInvitationControllerImpl) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	err = controller.Service.Delete(ctx, uint(id))
	if err != nil {
		if err == e.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrInvitationHasBeenProceed {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Team Invitation Success", nil)
}

func (controller *TeamInvitationControllerImpl) GetList(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
		return
	}

	eventID, _ := strconv.Atoi(ctx.Query("event"))
	filter := model.FilterInvitation{
		EventID: uint(eventID),
		Status:  ctx.Query("status"),
	}

	data, err := controller.Service.GetList(ctx, filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Team Invitation Success", data)
}

func (controller *TeamInvitationControllerImpl) GetDetail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetail(ctx, uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Invitation Success", data)
}

func (controller *TeamInvitationControllerImpl) GetDetail2(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetail2(ctx, uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Team Invitation Success", data)
}
