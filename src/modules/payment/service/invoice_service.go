package service

import (
	"be-sagara-hackathon/src/modules/payment/model"
	"be-sagara-hackathon/src/modules/payment/repository"
	"be-sagara-hackathon/src/utils"
)

type InvoiceService interface {
	GetList(
		filter model.FilterInvoice,
		pg *utils.PaginateQueryOffset,
	) (response model.ListInvoiceResponse, err error)
	GetDetail(id uint) (invoice model.InvoiceFull, err error)
	GetParticipantInvoice(participantID, eventID uint) (invoice model.InvoiceFull, err error)
}

type InvoiceServiceImpl struct {
	Repository repository.InvoiceRepository
}

func NewInvoiceService(repository repository.InvoiceRepository) InvoiceService {
	return &InvoiceServiceImpl{Repository: repository}
}

func (service *InvoiceServiceImpl) GetList(
	filter model.FilterInvoice,
	pg *utils.PaginateQueryOffset,
) (response model.ListInvoiceResponse, err error) {
	response.Invoices, response.TotalItem, response.TotalPage, err = service.Repository.FindAll(filter, pg)
	if err != nil {
		return
	}
	return
}

func (service *InvoiceServiceImpl) GetDetail(id uint) (invoice model.InvoiceFull, err error) {
	if invoice, err = service.Repository.FindOne(id); err != nil {
		return
	}
	return
}

func (service *InvoiceServiceImpl) GetParticipantInvoice(participantID, eventID uint) (invoice model.InvoiceFull, err error) {
	if invoice, err = service.Repository.FindByParticipantIDAndEventID(participantID, eventID); err != nil {
		return
	}
	return
}
