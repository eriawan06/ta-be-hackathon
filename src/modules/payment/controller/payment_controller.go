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

type PaymentController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	GetManyByInvoiceID(ctx *gin.Context)
}

type PaymentControllerImpl struct {
	Service service.PaymentService
}

func NewPaymentController(paymentService service.PaymentService) PaymentController {
	return &PaymentControllerImpl{Service: paymentService}
}

// Create Save New Payment godoc
// @Tags Payments
// @Summary Save New Payment
// @Description Save New Payment
// @Produce  json
// @Security ApiKeyAuth
// @Param body body dto.CreatePaymentRequest true "Body Request"
// @Success 201 {object} src.BaseSuccess
// @Success 400 {object} src.BaseFailure
// @Router /payments [post]
func (controller *PaymentControllerImpl) Create(ctx *gin.Context) {
	var request model.CreatePaymentRequest
	errorBinding := ctx.ShouldBindJSON(&request)
	if errorBinding != nil {
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
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		if err == e.ErrHaveUnprocessedPayment ||
			err == e.ErrInvoiceIsPaid {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusCreated, "Create Payment Success", nil)
}

// Update Update/Process Payment godoc
// @Tags Payments
// @Summary Update/Process Payment
// @Description Update/Process Payment. Admin Only
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Payment Id"
// @Param body body dto.UpdatePaymentRequest true "Body Request"
// @Success 201 {object} src.BaseSuccess
// @Success 400 {object} src.BaseFailure
// @Router /payments/{id} [put]
func (controller *PaymentControllerImpl) Update(ctx *gin.Context) {
	var request model.UpdatePaymentRequest
	errorBinding := ctx.ShouldBindJSON(&request)
	if errorBinding != nil {
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

	paymentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	if err = controller.Service.Update(ctx, uint(paymentID), request); err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Update Payment Success", nil)
}

// GetList Get List Payment godoc
// @Tags Payments
// @Summary Get All Payment
// @Description Get All Payment. Admin Only
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetListPaymentSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /payments [get]
func (controller *PaymentControllerImpl) GetList(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
		return
	}

	filter := model.FilterPayment{
		CreatedAt:   ctx.Query("created_at"),
		Status:      ctx.Query("status"),
		PaymentType: ctx.Query("type"),
		Search:      ctx.Query("q"),
	}

	data, err := controller.Service.GetList(filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Payment Success", data)
}

// GetDetail Get Detail Payment godoc
// @Tags Payments
// @Summary Get Detail Payment
// @Description Get Detail Payment
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Payment Id"
// @Success 200 {object} src.GetPaymentDetailSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /payments/{id} [get]
func (controller *PaymentControllerImpl) GetDetail(ctx *gin.Context) {
	paymentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetDetail(uint(paymentID))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Payment Success", data)
}

// GetManyByInvoiceID Get Payments By Invoice ID godoc
// @Tags Payments
// @Summary Get Payments By Invoice Id
// @Description Get Payments By Invoice Id
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetListPaymentSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /invoices/{id}/payments [get]
func (controller *PaymentControllerImpl) GetManyByInvoiceID(ctx *gin.Context) {
	invoiceID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Invoice Id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetManyByInvoiceID(uint(invoiceID))
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Payments By Invoice ID Success", data)
}
