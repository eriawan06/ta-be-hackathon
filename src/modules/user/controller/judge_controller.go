package controller

import (
	"be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/modules/user/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type JudgeController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
}

type JudgeControllerImpl struct {
	Service service.JudgeService
}

func NewJudgeController(judgeService service.JudgeService) JudgeController {
	return &JudgeControllerImpl{Service: judgeService}
}

// Create CreateJudge godoc
// @Tags User
// @Summary Create New Judge
// @Description Create New Judge
// @Produce  json
// @Security ApiKeyAuth
// @Param body body dto.CreateJudgeRequest true "Body Request"
// @Success 201 {object} src.BaseSuccess
// @Success 400 {object} src.BaseFailure
// @Router /users/judges [post]
func (controller *JudgeControllerImpl) Create(ctx *gin.Context) {
	var request model.CreateJudgeRequest
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

	err := controller.Service.Create(ctx, request)
	if err != nil {
		if err == e.ErrEmailAlreadyExists || err == e.ErrPhoneNumberAlreadyExists {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Judge Success", nil)
}

// Update UpdateJudge godoc
// @Tags User
// @Summary Update Judge
// @Description Update Judge
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Judge Id"
// @Param body body dto.UpdateJudgeRequest true "Body Request"
// @Success 200 {object} src.BaseSuccess
// @Success 400 {object} src.BaseFailure
// @Router /users/judges/{id} [put]
func (controller *JudgeControllerImpl) Update(ctx *gin.Context) {
	var request model.UpdateJudgeRequest
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

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	err = controller.Service.Update(ctx, request, uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Judge Success", nil)
}

// Delete DeleteJudge godoc
// @Tags User
// @Summary Delete Judge
// @Description Delete Judge
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Judge Id"
// @Success 200 {object} src.BaseSuccess
// @Success 400 {object} src.BaseFailure
// @Router /users/judges/{id} [delete]
func (controller *JudgeControllerImpl) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	err = controller.Service.Delete(ctx, uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Judge Success", nil)
}

// GetList Get List Judges godoc
// @Tags User
// @Summary Get List Judges
// @Description Get List Judges
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetListJudgeSuccess
// @Success 400 {object} src.BaseFailure
// @Router /users/judges [get]
func (controller *JudgeControllerImpl) GetList(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	var filter model.FilterUser
	filter.Status = ctx.Query("status")
	filter.Search = ctx.Query("q")

	data, err := controller.Service.GetList(filter, pg)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Judge Success", data)
}

// GetDetail Get Judge By Id godoc
// @Tags User
// @Summary Get Judge By Id
// @Description Get Judge By Id
// @Param id path int true "Judge Id"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetJudgeSuccess
// @Success 400 {object} src.BaseFailure
// @Router /users/judges/{id} [get]
func (controller *JudgeControllerImpl) GetDetail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetail(uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Judge Success", data)
}
