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

type ProjectAssessmentController interface {
	Create(ctx *gin.Context)
	GetByProjectID(ctx *gin.Context)
	GetByJudgeAndProjectID(ctx *gin.Context)
}

type ProjectAssessmentControllerImpl struct {
	Service service.ProjectAssessmentService
}

func NewProjectAssessmentController(assessmentService service.ProjectAssessmentService) ProjectAssessmentController {
	return &ProjectAssessmentControllerImpl{Service: assessmentService}
}

func (controller *ProjectAssessmentControllerImpl) Create(ctx *gin.Context) {
	var request model.CreateBatchProjectAssessmentRequest
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

	projectID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid project id", []string{err.Error()})
		return
	}

	err = controller.Service.CreateBatch(ctx, uint(projectID), request)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		if err == e.ErrProjectStatusShouldBeSubmitted {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Project Assessment Success", nil)
}

func (controller *ProjectAssessmentControllerImpl) GetByProjectID(ctx *gin.Context) {
	projectID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid project id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetByProjectID(uint(projectID))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Project Assessment Success", data)
}

func (controller *ProjectAssessmentControllerImpl) GetByJudgeAndProjectID(ctx *gin.Context) {
	projectID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid project id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetByJudgeAndProjectID(ctx, uint(projectID))
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

	common.SendSuccess(ctx, http.StatusOK, "Get Project Assessment By Judge Success", data)
}
