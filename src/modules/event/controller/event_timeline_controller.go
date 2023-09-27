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

type EventTimelineController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
}

type EventTimelineControllerImpl struct {
	Service service.EventTimelineService
}

func NewEventTimelineController(service service.EventTimelineService) EventTimelineController {
	return &EventTimelineControllerImpl{Service: service}
}

func (controller *EventTimelineControllerImpl) Create(ctx *gin.Context) {
	var request model.EventTimelineRequest
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

	common.SendSuccess(ctx, http.StatusCreated, "Create Event Timeline Success", nil)
}

func (controller *EventTimelineControllerImpl) Update(ctx *gin.Context) {
	var request model.EventTimelineRequest
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

	etlID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event timeline id", []string{err.Error()})
		return
	}

	if err = controller.Service.Update(ctx, request, uint(etlID)); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Event Timeline Success", nil)
}

func (controller *EventTimelineControllerImpl) Delete(ctx *gin.Context) {
	etlID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event timeline id", []string{err.Error()})
		return
	}

	if err = controller.Service.Delete(uint(etlID)); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Event Timeline Success", nil)
}

func (controller *EventTimelineControllerImpl) GetList(ctx *gin.Context) {
	eventID, err := strconv.Atoi(ctx.Query("event"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{"Invalid event"})
		return
	}
	filter := model.FilterEventTimeline{
		EventID: uint(eventID),
	}
	data, err := controller.Service.GetList(filter)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Event Timeline Success", data)
}

func (controller *EventTimelineControllerImpl) GetDetail(ctx *gin.Context) {
	etlID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event timeline id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetail(uint(etlID))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Event Timeline Success", data)
}
