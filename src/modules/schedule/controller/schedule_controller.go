package controller

import (
	"be-sagara-hackathon/src/modules/schedule/model"
	"be-sagara-hackathon/src/modules/schedule/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ScheduleController interface {
	CreateSchedule(ctx *gin.Context)
	CreateScheduleTeam(ctx *gin.Context)
	UpdateSchedule(ctx *gin.Context)
	DeleteSchedule(ctx *gin.Context)
	DeleteScheduleTeam(ctx *gin.Context)
	GetListSchedule(ctx *gin.Context)
	GetDetailSchedule(ctx *gin.Context)
}

type ScheduleControllerImpl struct {
	Service service.ScheduleService
}

func NewScheduleController(service service.ScheduleService) ScheduleController {
	return &ScheduleControllerImpl{Service: service}
}

func (controller *ScheduleControllerImpl) CreateSchedule(ctx *gin.Context) {
	var request model.ScheduleRequest
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

	data, err := controller.Service.CreateSchedule(ctx, request)
	if err != nil {
		if err == e.ErrEventNotRunning || err == e.ErrScheduleDateNotValid {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Schedule Success", data)
}

func (controller *ScheduleControllerImpl) CreateScheduleTeam(ctx *gin.Context) {
	var request model.ScheduleTeam
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

	err := controller.Service.CreateScheduleTeam(request)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Schedule Team Success", nil)
}

func (controller *ScheduleControllerImpl) UpdateSchedule(ctx *gin.Context) {
	var request model.ScheduleRequest
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

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	err = controller.Service.UpdateSchedule(ctx, uint(id), request)
	if err != nil {
		if err == e.ErrScheduleDateNotValid {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request Error", []string{err.Error()})
			return
		}

		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Schedule Success", nil)
}

func (controller *ScheduleControllerImpl) DeleteSchedule(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	err = controller.Service.DeleteSchedule(uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Schedule Success", nil)
}

func (controller *ScheduleControllerImpl) DeleteScheduleTeam(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid schedule Id", []string{err.Error()})
		return
	}

	teamID, err := strconv.Atoi(ctx.Param("team_id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid team Id", []string{err.Error()})
		return
	}

	err = controller.Service.DeleteScheduleTeam(uint(id), uint(teamID))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Schedule Team Success", nil)
}

func (controller *ScheduleControllerImpl) GetListSchedule(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	eventID, _ := strconv.Atoi(ctx.Query("event"))
	filter := model.FilterSchedule{
		EventID: uint(eventID),
		HeldOn:  ctx.Query("held"),
		Search:  ctx.Query("q"),
	}

	data, err := controller.Service.GetListSchedule(ctx, filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Schedule Success", data)
}

func (controller *ScheduleControllerImpl) GetDetailSchedule(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetailSchedule(uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Schedule Success", data)
}
