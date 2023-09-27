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

type EventRuleController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	GetActiveByEventID(ctx *gin.Context)
}

type EventRuleControllerImpl struct {
	Service service.EventRuleService
}

func NewEventRuleController(service service.EventRuleService) EventRuleController {
	return &EventRuleControllerImpl{Service: service}
}

func (controller *EventRuleControllerImpl) Create(ctx *gin.Context) {
	var request model.EventRuleRequest
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

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Event Rule Success", nil)
}

func (controller *EventRuleControllerImpl) Update(ctx *gin.Context) {
	var request model.UpdateEventRuleRequest
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

	erID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event company id", []string{err.Error()})
		return
	}

	if err = controller.Service.Update(ctx, request, uint(erID)); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Event Rule Success", nil)
}

func (controller *EventRuleControllerImpl) Delete(ctx *gin.Context) {
	erID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event company id", []string{err.Error()})
		return
	}

	if err = controller.Service.Delete(uint(erID)); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Event Rule Success", nil)
}

func (controller *EventRuleControllerImpl) GetList(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	eventID, err := strconv.Atoi(ctx.Query("event"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{"Invalid event"})
		return
	}
	filter := model.FilterEventRule{
		EventID: uint(eventID),
		Status:  ctx.Query("status"),
	}
	data, err := controller.Service.GetList(filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Event Rule Success", data)
}

func (controller *EventRuleControllerImpl) GetDetail(ctx *gin.Context) {
	erID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetail(uint(erID))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendSuccess(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal server error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Event Rule Success", data)
}

func (controller *EventRuleControllerImpl) GetActiveByEventID(ctx *gin.Context) {
	eventID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetActiveByEventID(uint(eventID))
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal server error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Event Rules Success", data)
}
