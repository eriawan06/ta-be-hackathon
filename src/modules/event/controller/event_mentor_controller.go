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

type EventMentorController interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
}

type EventMentorControllerImpl struct {
	Service service.EventMentorService
}

func NewEventMentorController(service service.EventMentorService) EventMentorController {
	return &EventMentorControllerImpl{Service: service}
}

// Create Create Event Mentor godoc
// @Tags Events
// @Summary Create Event Mentor
// @Description Create Event Mentor
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.CreateEventMentorRequest true "Body Request"
// @Success 201 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events/mentors [post]
func (controller *EventMentorControllerImpl) Create(ctx *gin.Context) {
	var request model.EventMentorRequest
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
		if err == e.ErrEventMentorExist {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Event Mentor Success", nil)
}

// Delete Delete Event Mentor godoc
// @Tags Events
// @Summary Delete Event Mentor
// @Description Delete Event Mentor
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Event Mentor ID"
// @Success 200 {object} src.BaseSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events/mentors/{id} [delete]
func (controller *EventMentorControllerImpl) Delete(ctx *gin.Context) {
	eventMentorID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event mentor id", []string{err.Error()})
		return
	}

	if err = controller.Service.Delete(uint(eventMentorID)); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Event Mentor Success", nil)
}

// GetAll Get All Event Mentor godoc
// @Tags Events
// @Summary Get All Event Mentor
// @Description Get All Event Mentor
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetListEventMentorSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events/mentors [get]
func (controller *EventMentorControllerImpl) GetAll(ctx *gin.Context) {
	eventID, err := strconv.Atoi(ctx.Query("event"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{"Invalid event"})
		return
	}
	filter := model.FilterEventMentor{
		EventID: uint(eventID),
	}

	data, err := controller.Service.GetAll(filter)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get All Event Mentor Success", data)
}

// GetDetail Get One Event Mentor godoc
// @Tags Events
// @Summary Get One Event Mentor By Id
// @Description Get One Event Mentor By Id
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Event Mentor ID"
// @Success 200 {object} src.GetEventMentorSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /events/mentors/{id} [get]
func (controller *EventMentorControllerImpl) GetDetail(ctx *gin.Context) {
	eventMentorID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event mentor id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetail(uint(eventMentorID))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Event Mentor Success", data)
}
