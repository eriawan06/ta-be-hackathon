package controller

import (
	"be-sagara-hackathon/src/modules/project/model"
	"be-sagara-hackathon/src/modules/project/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProjectController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	UpdateStatus(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}

type ProjectControllerImpl struct {
	Service service.ProjectService
}

func NewProjectController(projectService service.ProjectService) ProjectController {
	return &ProjectControllerImpl{Service: projectService}
}

func (controller *ProjectControllerImpl) Create(ctx *gin.Context) {
	var request model.CreateProjectRequest
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
	request.Action = "create"
	if errs := utils.NewCustomValidator().ValidateStruct(request); errs != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request", errs)
		return
	}

	err = controller.Service.Create(ctx, request)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		if err == e.ErrEventNotRunning || err == e.ErrTeamAlreadyHasProject {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Project Success", nil)
}

func (controller *ProjectControllerImpl) Update(ctx *gin.Context) {
	var request model.UpdateProjectRequest
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
	request.Action = "update"
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
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		if err == e.ErrEventNotRunning || err == e.ErrProjectStatusShouldBeDraft {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Project Success", nil)
}

func (controller *ProjectControllerImpl) UpdateStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	status := ctx.Param("status")
	if status == "" {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{e.ErrInvalidStatus.Error()})
		return
	}

	err = controller.Service.UpdateStatus(ctx, uint(id), status)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrEventNotRunning || err == e.ErrInvalidStatus {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Status Project Success", nil)
}

func (controller *ProjectControllerImpl) GetDetail(ctx *gin.Context) {
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

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Project Success", data)
}

func (controller *ProjectControllerImpl) GetAll(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
		return
	}

	eventID, _ := strconv.Atoi(ctx.Query("event"))
	teamID, _ := strconv.Atoi(ctx.Query("team"))
	filter := model.FilterProject{
		EventID:   uint(eventID),
		TeamID:    uint(teamID),
		CreatedAt: ctx.Query("created_at"),
		Status:    ctx.Query("status"),
		Search:    ctx.Query("q"),
	}

	data, err := controller.Service.GetAll(filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get All Project Success", data)
}
