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

type EventController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	GetLatest(ctx *gin.Context)
	GetSchedules(ctx *gin.Context)
}

type EventControllerImpl struct {
	Service service.EventService
}

func NewEventController(service service.EventService) EventController {
	return &EventControllerImpl{Service: service}
}

// Create Create Event godoc
// @Tags Events
// @Summary Create Event
// @Description Create Event
// @Produce  json
// @Security ApiKeyAuth
// @Param body body model.CreateEventRequest true "Body Request"
// @Success 200 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events [post]
func (controller *EventControllerImpl) Create(ctx *gin.Context) {
	var request model.CreateEventRequest
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

	data, err := controller.Service.CreateEvent(ctx, request)
	if err != nil {
		if err == e.ErrLatestEventRunning {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Create Event Success", data)
}

// Update Update Event godoc
// @Tags Events
// @Summary Update Event
// @Description Update Event
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Event ID"
// @Param body body model.UpdateEventRequest true "Body Request"
// @Success 200 {object} src.BaseSuccess
// @Success 400 {object} src.BaseFailure
// @Router /events/{id} [put]
func (controller *EventControllerImpl) Update(ctx *gin.Context) {
	var request model.UpdateEventRequest
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

	eventID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event id", []string{err.Error()})
		return
	}

	if err = controller.Service.UpdateEvent(ctx, request, uint(eventID)); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		if err == e.ErrLatestEventRunning {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Event Success", nil)
}

// Delete Delete Event godoc
// @Tags Events
// @Summary Delete Event
// @Description Delete Event
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Event ID"
// @Success 200 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events/{id} [delete]
func (controller *EventControllerImpl) Delete(ctx *gin.Context) {
	eventID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event id", []string{err.Error()})
		return
	}

	if err = controller.Service.DeleteEvent(ctx, uint(eventID)); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Event Success", nil)
}

// GetList Get List Event godoc
// @Tags Events
// @Summary Get All Events
// @Description Get All Events
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetListEventsSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events [get]
func (controller *EventControllerImpl) GetList(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	filter := model.FilterEvent{
		Search:    ctx.Query("q"),
		Status:    ctx.Query("status"),
		StartDate: ctx.Query("start"),
		EndDate:   ctx.Query("end"),
	}
	data, err := controller.Service.GetListEvent(filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Event Success", data)
}

// GetDetail Get Detail Event By Event ID godoc
// @Tags Events
// @Summary Get Event By Event ID
// @Description Get Event By Event ID
// @Param id path int true "Event ID"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetEventSuccess
// @Success 400 {object} src.BaseFailure
// @Router /events/{id} [get]
func (controller *EventControllerImpl) GetDetail(ctx *gin.Context) {
	eventId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetailEvent(uint(eventId))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendSuccess(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal server error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Event Success", data)
}

// GetLatest Get Latest Event godoc
// @Tags Events
// @Summary Get Latest Event
// @Description Get Latest Event
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetEventSuccess
// @Success 400 {object} src.BaseFailure
// @Router /events/latest [get]
func (controller *EventControllerImpl) GetLatest(ctx *gin.Context) {
	data, err := controller.Service.GetLatestEvent()
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendSuccess(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Latest Event Success", data)
}

func (controller *EventControllerImpl) GetSchedules(ctx *gin.Context) {
	eventID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetSchedules(ctx, uint(eventID))
	if err != nil {
		if err == e.ErrForbidden {
			common.SendSuccess(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		if err == e.ErrDataNotFound {
			common.SendSuccess(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Event Schedule Success", data)
}
