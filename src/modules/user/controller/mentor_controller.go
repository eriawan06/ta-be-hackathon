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

type MentorController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
}

type MentorControllerImpl struct {
	Service service.MentorService
}

func NewMentorController(mentorService service.MentorService) MentorController {
	return &MentorControllerImpl{Service: mentorService}
}

// Create CreateMentor godoc
// @Tags User
// @Summary Create New Mentor
// @Description Create New Mentor
// @Produce  json
// @Security ApiKeyAuth
// @Param body body dto.CreateMentorRequest true "Body Request"
// @Success 201 {object} src.BaseSuccess
// @Success 400 {object} src.BaseFailure
// @Router /users/mentors [post]
func (controller *MentorControllerImpl) Create(ctx *gin.Context) {
	var request model.CreateMentorRequest
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

	common.SendSuccess(ctx, http.StatusCreated, "Create Mentor Success", nil)
}

// Update UpdateMentor godoc
// @Tags User
// @Summary Update Mentor
// @Description Update Mentor
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Mentor Id"
// @Param body body dto.UpdateMentorRequest true "Body Request"
// @Success 200 {object} src.BaseSuccess
// @Success 400 {object} src.BaseFailure
// @Router /users/mentors/{id} [put]
func (controller *MentorControllerImpl) Update(ctx *gin.Context) {
	var request model.UpdateMentorRequest
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

	common.SendSuccess(ctx, http.StatusOK, "Update Mentor Success", nil)
}

// Delete DeleteMentor godoc
// @Tags User
// @Summary Delete Mentor
// @Description Delete Mentor
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Mentor Id"
// @Success 200 {object} src.BaseSuccess
// @Success 400 {object} src.BaseFailure
// @Router /users/mentors/{id} [delete]
func (controller *MentorControllerImpl) Delete(ctx *gin.Context) {
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

	common.SendSuccess(ctx, http.StatusOK, "Delete Mentor Success", nil)
}

// GetList Get List Mentors godoc
// @Tags User
// @Summary Get List Mentors
// @Description Get List Mentors
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetListMentorSuccess
// @Success 400 {object} src.BaseFailure
// @Router /users/mentors [get]
func (controller *MentorControllerImpl) GetList(ctx *gin.Context) {
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

	common.SendSuccess(ctx, http.StatusOK, "Get List Mentor Success", data)
}

// GetDetail Get Mentor By Id godoc
// @Tags User
// @Summary Get Mentor By Id
// @Description Get Mentor By Id
// @Param id path int true "Mentor Id"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetMentorSuccess
// @Success 400 {object} src.BaseFailure
// @Router /users/mentors/{id} [get]
func (controller *MentorControllerImpl) GetDetail(ctx *gin.Context) {
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

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Mentor Success", data)
}
