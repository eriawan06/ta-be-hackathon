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

type EventJudgeController interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
}

type EventJudgeControllerImpl struct {
	Service service.EventJudgeService
}

func NewEventJudgeController(service service.EventJudgeService) EventJudgeController {
	return &EventJudgeControllerImpl{Service: service}
}

// Create Create Event Judge godoc
// @Tags Events
// @Summary Create Event Judge
// @Description Create Event Judge
// @Produce  json
// @Security ApiKeyAuth
// @Param body body dto.CreateEventJudgeRequest true "Body Request"
// @Success 201 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events/judges [post]
func (controller *EventJudgeControllerImpl) Create(ctx *gin.Context) {
	var request model.EventJudgeRequest
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

	if err := controller.Service.Create(ctx, request); err != nil {
		if err == e.ErrEventJudgeExist {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Event Judge Success", nil)
}

// Delete Delete Event Judge godoc
// @Tags Events
// @Summary Delete Event Judge
// @Description Delete Event Judge
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Event Judge ID"
// @Success 200 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events/judges/{id} [delete]
func (controller *EventJudgeControllerImpl) Delete(ctx *gin.Context) {
	eventJudgeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event judge id", []string{err.Error()})
		return
	}

	if err = controller.Service.Delete(uint(eventJudgeID)); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Event Judge Success", nil)
}

// GetAll Get All Event Judge godoc
// @Tags Events
// @Summary Get All Event Judge
// @Description Get All Event Judge
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetListEventJudgeSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events/judges [get]
func (controller *EventJudgeControllerImpl) GetAll(ctx *gin.Context) {
	eventID, err := strconv.Atoi(ctx.Query("event"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{"Invalid event"})
		return
	}
	filter := model.FilterEventJudge{
		EventID: uint(eventID),
	}

	data, err := controller.Service.GetAll(filter)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get All Event Judge Success", data)
}

// GetOne Get One Event Judge godoc
// @Tags Events
// @Summary Get One Event Judge By Id
// @Description Get One Event Judge By Id
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Event Judge ID"
// @Success 200 {object} src.GetEventJudgeSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events/judges/{id} [get]
func (controller *EventJudgeControllerImpl) GetDetail(ctx *gin.Context) {
	eventJudgeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event mentor id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetail(uint(eventJudgeID))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Event Judge Success", data)
}
