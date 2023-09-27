package controller

import (
	"be-sagara-hackathon/src/modules/payment/model"
	"be-sagara-hackathon/src/modules/payment/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PaymentMethodController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetOne(ctx *gin.Context)
}

type PaymentMethodControllerImpl struct {
	Service service.PaymentMethodService
}

func NewPaymentMethodController(service service.PaymentMethodService) PaymentMethodController {
	return &PaymentMethodControllerImpl{Service: service}
}

func (controller *PaymentMethodControllerImpl) Create(ctx *gin.Context) {
	var request model.PaymentMethodRequest
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

	err := controller.Service.Create(ctx, request)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Payment Method Success", nil)
}

func (controller *PaymentMethodControllerImpl) Update(ctx *gin.Context) {
	var request model.UpdatePaymentMethodRequest
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

	err = controller.Service.Update(ctx, request, uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Payment Method Success", nil)
}

func (controller *PaymentMethodControllerImpl) Delete(ctx *gin.Context) {
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

	common.SendSuccess(ctx, http.StatusOK, "Delete Payment Method Success", nil)
}

func (controller *PaymentMethodControllerImpl) GetList(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
		return
	}

	var filter model.FilterPaymentMethod
	filter.Status = ctx.Query("status")

	data, err := controller.Service.GetList(filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Payment Method Success", data)
}

func (controller *PaymentMethodControllerImpl) GetOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetOne(uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get One Payment Method Success", data)
}
