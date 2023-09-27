package controller

import (
	"be-sagara-hackathon/src/modules/master-data/occupation/model"
	"be-sagara-hackathon/src/modules/master-data/occupation/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OccupationController interface {
	CreateOccupation(ctx *gin.Context)
	UpdateOccupation(ctx *gin.Context)
	GetListOccupation(ctx *gin.Context)
	GetDetailOccupation(ctx *gin.Context)
}

type OccupationControllerImpl struct {
	Service service.OccupationService
}

func NewOccupationController(service service.OccupationService) OccupationController {
	return &OccupationControllerImpl{Service: service}
}

func (controller *OccupationControllerImpl) CreateOccupation(ctx *gin.Context) {
	var request model.OccupationRequest
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

	err := controller.Service.CreateOccupation(ctx, request)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	// Return Success Login
	common.SendSuccess(ctx, http.StatusCreated, "Create Occupation Success", nil)
}

func (controller *OccupationControllerImpl) UpdateOccupation(ctx *gin.Context) {
	var request model.UpdateOccupationRequest
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

	err = controller.Service.UpdateOccupation(ctx, request, uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Occupation Success", nil)
}

func (controller *OccupationControllerImpl) GetListOccupation(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	var filter model.FilterOccupation
	filter.Name = ctx.Query("name")
	filter.Status = ctx.Query("status")

	data, err := controller.Service.GetListOccupation(filter, pg)
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Occupation Success", data)
}

func (controller *OccupationControllerImpl) GetDetailOccupation(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetailOccupation(uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Occupation Success", data)
}
