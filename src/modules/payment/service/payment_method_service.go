package service

import (
	"be-sagara-hackathon/src/modules/payment/model"
	"be-sagara-hackathon/src/modules/payment/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"context"
)

type PaymentMethodService interface {
	Create(ctx context.Context, request model.PaymentMethodRequest) error
	Update(ctx context.Context, request model.UpdatePaymentMethodRequest, id uint) error
	Delete(ctx context.Context, id uint) error
	GetList(
		filter model.FilterPaymentMethod,
		pg *utils.PaginateQueryOffset,
	) (response model.ListPaymentMethodResponse, err error)
	GetOne(id uint) (method model.PaymentMethod, err error)
}

type PaymentMethodServiceImpl struct {
	Repository repository.PaymentMethodRepository
}

func NewPaymentMethod(repository repository.PaymentMethodRepository) PaymentMethodService {
	return &PaymentMethodServiceImpl{Repository: repository}
}

func (service *PaymentMethodServiceImpl) Create(ctx context.Context, request model.PaymentMethodRequest) error {
	if err := service.Repository.Save(model.PaymentMethod{
		BaseEntity:    builder.BuildBaseEntity(ctx, true, nil),
		Name:          request.Name,
		BankCode:      request.BankCode,
		AccountNumber: request.AccountNumber,
		AccountName:   request.AccountName,
		IsActive:      true,
	}); err != nil {
		return err
	}
	return nil
}

func (service *PaymentMethodServiceImpl) Update(ctx context.Context, request model.UpdatePaymentMethodRequest, id uint) error {
	method, err := service.Repository.FindByID(id)
	if err != nil {
		return err
	}

	if err = service.Repository.Update(id, model.PaymentMethod{
		BaseEntity:    builder.BuildBaseEntity(ctx, false, &method.BaseEntity),
		Name:          request.Name,
		BankCode:      request.BankCode,
		AccountNumber: request.AccountNumber,
		AccountName:   request.AccountName,
		IsActive:      request.IsActive,
	}); err != nil {
		return err
	}
	return nil
}

func (service *PaymentMethodServiceImpl) Delete(ctx context.Context, id uint) error {
	authenticatedUser := ctx.Value("user").(um.User)
	if _, err := service.Repository.FindByID(id); err != nil {
		return err
	}

	if err := service.Repository.Delete(id, authenticatedUser.Email); err != nil {
		return err
	}

	return nil
}

func (service *PaymentMethodServiceImpl) GetList(filter model.FilterPaymentMethod, pg *utils.PaginateQueryOffset) (response model.ListPaymentMethodResponse, err error) {
	response.Methods, response.TotalItem, response.TotalPage, err = service.Repository.FindAll(filter, pg)
	if err != nil {
		return
	}
	return
}

func (service *PaymentMethodServiceImpl) GetOne(id uint) (method model.PaymentMethod, err error) {
	method, err = service.Repository.FindByID(id)
	if err != nil {
		return
	}
	return
}
