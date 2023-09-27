package controller

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/modules/event/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type EventCompanyController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
}

type EventCompanyControllerImpl struct {
	Service service.EventCompanyService
}

func NewEventCompanyController(service service.EventCompanyService) EventCompanyController {
	return &EventCompanyControllerImpl{Service: service}
}

func (controller *EventCompanyControllerImpl) Create(ctx *gin.Context) {
	var request model.EventCompanyRequest
	if errorBinding := ctx.ShouldBindJSON(&request); errorBinding != nil {
		if errorBinding.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(errorBinding))
		return
	}

	// Validate request body
	request.Action = "create"
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	if err := controller.Service.Create(ctx, request); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		if err == e.ErrNameAlreadyExists || err == e.ErrPhoneNumberAlreadyExists || err == e.ErrEmailAlreadyExists {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Event Company Success", nil)
}

func (controller *EventCompanyControllerImpl) Update(ctx *gin.Context) {
	var request model.EventCompanyRequest
	if errorBinding := ctx.ShouldBindJSON(&request); errorBinding != nil {
		if errorBinding.Error() == "EOF" {
			common.SendError(ctx, http.StatusBadRequest, "Body is empty", []string{"Body required"})
			return
		}

		common.SendError(ctx, http.StatusBadRequest, "Invalid request", utils.SplitError(errorBinding))
		return
	}

	// Validate request body
	request.Action = "update"
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	ecID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event company id", []string{err.Error()})
		return
	}

	if err = controller.Service.Update(ctx, request, uint(ecID)); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		if err == e.ErrNameAlreadyExists || err == e.ErrPhoneNumberAlreadyExists || err == e.ErrEmailAlreadyExists {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Event Company Success", nil)
}

func (controller *EventCompanyControllerImpl) Delete(ctx *gin.Context) {
	ecID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event company id", []string{err.Error()})
		return
	}

	if err = controller.Service.Delete(uint(ecID)); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Event Company Success", nil)
}

func (controller *EventCompanyControllerImpl) GetList(ctx *gin.Context) {
	eventID, err := strconv.Atoi(ctx.Query("event"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{"Invalid event"})
		return
	}
	filter := model.FilterEventCompany{
		EventID: uint(eventID),
	}
	data, err := controller.Service.GetList(filter)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Event Company Success", data)
}

func (controller *EventCompanyControllerImpl) GetDetail(ctx *gin.Context) {
	ecID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event company id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetail(uint(ecID))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendSuccess(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal server error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Event Company Success", data)
}
