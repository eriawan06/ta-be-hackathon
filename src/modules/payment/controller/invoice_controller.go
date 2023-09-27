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

type InvoiceController interface {
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	GetParticipantInvoice(ctx *gin.Context)
}

type InvoiceControllerImpl struct {
	Service service.InvoiceService
}

func NewInvoiceController(service service.InvoiceService) InvoiceController {
	return &InvoiceControllerImpl{Service: service}
}

// GetList Get List Invoice godoc
// @Tags Payments
// @Summary Get All Invoices
// @Description Get All Invoices
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} src.GetListInvoicesFullSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /invoices [get]
func (controller *InvoiceControllerImpl) GetList(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
		return
	}

	eventID, _ := strconv.Atoi(ctx.Query("event"))
	filter := model.FilterInvoice{
		Search:  ctx.Query("q"),
		EventID: uint(eventID),
		Status:  ctx.Query("status"),
	}

	data, err := controller.Service.GetList(filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Invoice Success", data)
}

// GetDetail Get Detail Invoice godoc
// @Tags Payments
// @Summary Get Invoice
// @Description Get Invoice
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Invoice Id"
// @Success 200 {object} src.GetInvoiceFullSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /invoices/{id} [get]
func (controller *InvoiceControllerImpl) GetDetail(ctx *gin.Context) {
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

	common.SendSuccess(ctx, http.StatusOK, "Get Detail Invoice Success", data)
}

// GetParticipantInvoice Get Participant's Invoice godoc
// @Tags Participants
// @Summary Get Participant Invoice
// @Description Get Participant Invoice
// @Produce json
// @Security ApiKeyAuth
// @Param participant_id path int true "Participant Id"
// @Param event_id path int true "Event Id"
// @Success 200 {object} src.GetInvoiceFullSuccess
// @Failure 400 {object} src.BaseFailure
// @Router /users/participants/{participant_id}/events/{event_id}/invoice [get]
func (controller *InvoiceControllerImpl) GetParticipantInvoice(ctx *gin.Context) {
	participantID, err := strconv.Atoi(ctx.Param("participant_id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid participant id", []string{err.Error()})
		return
	}

	eventID, err := strconv.Atoi(ctx.Param("event_id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid event id", []string{err.Error()})
		return
	}

	data, err := controller.Service.GetParticipantInvoice(uint(participantID), uint(eventID))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found Error", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Participant's Invoice Success", data)
}
