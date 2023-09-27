package service

import (
	"be-sagara-hackathon/src/modules/payment/model"
	"be-sagara-hackathon/src/modules/payment/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	ur "be-sagara-hackathon/src/modules/user/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
	"context"
	"time"
)

type PaymentService interface {
	Create(ctx context.Context, request model.CreatePaymentRequest) error
	Update(ctx context.Context, paymentID uint, request model.UpdatePaymentRequest) error
	GetList(
		filter model.FilterPayment,
		pg *utils.PaginateQueryOffset,
	) (response model.ListPaymentResponse, err error)
	GetDetail(paymentID uint) (payment model.PaymentDetail, err error)
	GetManyByInvoiceID(invID uint) (payments []model.PaymentLite, err error)
}

type PaymentServiceImpl struct {
	Repository      repository.PaymentRepository
	InvoiceRepo     repository.InvoiceRepository
	ParticipantRepo ur.ParticipantRepository
}

func NewPaymentService(
	paymentRepository repository.PaymentRepository,
	invoiceRepository repository.InvoiceRepository,
	participantRepository ur.ParticipantRepository,
) PaymentService {
	return &PaymentServiceImpl{
		Repository:      paymentRepository,
		InvoiceRepo:     invoiceRepository,
		ParticipantRepo: participantRepository,
	}
}

func (service *PaymentServiceImpl) Create(ctx context.Context, request model.CreatePaymentRequest) error {
	//check unprocessed payment
	unprocessed, err := service.Repository.FindUnprocessedByInvoiceID(request.InvoiceID)
	if err != nil && err != e.ErrDataNotFound {
		return err
	}

	if unprocessed.ID != 0 {
		return e.ErrHaveUnprocessedPayment
	}

	// get related invoice
	invoice, err := service.InvoiceRepo.FindOne(request.InvoiceID)
	if err != nil {
		return err
	}

	//check invoice status
	if invoice.Status == constants.InvoicePaid {
		return e.ErrInvoiceIsPaid
	}

	//get related participant
	participant, err := service.ParticipantRepo.FindByID(invoice.ParticipantID)
	if err != nil {
		return err
	}

	if err = service.Repository.Save(model.Payment{
		BaseEntity:        builder.BuildBaseEntity(ctx, true, nil),
		InvoiceID:         request.InvoiceID,
		Invoice:           model.Invoice{ParticipantID: participant.ID},
		PaymentType:       constants.PaymentTypeManual,
		PaymentMethodID:   &request.PaymentMethodID,
		BankName:          &request.BankName,
		BankAccountName:   &request.AccountName,
		BankAccountNumber: &request.AccountNumber,
		Evidence:          &request.Evidence,
	}); err != nil {
		return err
	}
	return err
}

func (service *PaymentServiceImpl) Update(ctx context.Context, paymentID uint, request model.UpdatePaymentRequest) error {
	authenticatedUser := ctx.Value("user").(um.User)

	//check payment
	checkPayment, err := service.Repository.FindOne(paymentID)
	if err != nil {
		return err
	}

	//get related participant
	if _, err = service.ParticipantRepo.FindByID(checkPayment.ParticipantID); err != nil {
		return err
	}

	//get related invoice
	relatedInvoice, err := service.InvoiceRepo.FindOne(checkPayment.InvoiceID)
	if err != nil {
		return err
	}

	// Update Invoice
	amount := request.Amount - checkPayment.Amount
	paidAmount := relatedInvoice.PaidAmount + amount
	invUpdate := model.Invoice{PaidAmount: paidAmount, Status: constants.InvoiceProcessing}
	if invUpdate.PaidAmount >= relatedInvoice.Amount {
		invUpdate.Status = constants.InvoicePaid
		invUpdate.ApprovedAt = helper.ReferTime(time.Now())
		invUpdate.ApprovedBy = &authenticatedUser.Email
	}

	if relatedInvoice.Status == constants.InvoicePaid &&
		invUpdate.PaidAmount != relatedInvoice.Amount {
		invUpdate.Status = constants.InvoiceProcessing
		invUpdate.ApprovedAt = nil
		invUpdate.ApprovedBy = nil
	}

	invUpdate.ParticipantID = checkPayment.ParticipantID
	if err = service.Repository.Update(paymentID, model.Payment{
		BaseEntity: common.BaseEntity{
			UpdatedAt: time.Now(),
			UpdatedBy: authenticatedUser.Email,
		},
		Status:    constants.PaymentStatusProceed,
		Amount:    request.Amount,
		Note:      request.Note,
		ProceedAt: helper.ReferTime(time.Now()),
		ProceedBy: &authenticatedUser.Email,
		InvoiceID: checkPayment.InvoiceID,
		Invoice:   invUpdate,
	}); err != nil {
		return err
	}

	return nil
}

func (service *PaymentServiceImpl) GetList(
	filter model.FilterPayment,
	pg *utils.PaginateQueryOffset,
) (response model.ListPaymentResponse, err error) {
	response.Payments, response.TotalItem, response.TotalPage, err = service.Repository.FindAll(filter, pg)
	if err != nil {
		return
	}
	return
}

func (service *PaymentServiceImpl) GetDetail(paymentID uint) (payment model.PaymentDetail, err error) {
	payment, err = service.Repository.FindOne(paymentID)
	if err != nil {
		return
	}
	return
}

func (service *PaymentServiceImpl) GetManyByInvoiceID(invID uint) (payments []model.PaymentLite, err error) {
	if payments, err = service.Repository.FindManyByInvoiceID(invID); err != nil {
		return
	}
	return
}
