package service

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/utils/common/builder"
	"context"
)

type EventCompanyService interface {
	Create(ctx context.Context, req model.EventCompanyRequest) error
	Update(ctx context.Context, req model.EventCompanyRequest, ecID uint) error
	Delete(ecID uint) error
	GetList(filter model.FilterEventCompany) (companies []model.EventCompany, err error)
	GetDetail(ecID uint) (company model.EventCompany, err error)
}

type EventCompanyServiceImpl struct {
	Repository      repository.EventCompanyRepository
	EventRepository repository.EventRepository
}

func NewEventCompanyService(
	repository repository.EventCompanyRepository,
	eventRepository repository.EventRepository,
) EventCompanyService {
	return &EventCompanyServiceImpl{Repository: repository, EventRepository: eventRepository}
}

func (service *EventCompanyServiceImpl) Create(ctx context.Context, req model.EventCompanyRequest) error {
	if _, err := service.EventRepository.FindOne(req.EventID); err != nil {
		return err
	}

	if err := service.Repository.Save(model.EventCompany{
		BaseEntity:        builder.BuildBaseEntity(ctx, true, nil),
		EventID:           req.EventID,
		Name:              req.Name,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		PartnershipType:   req.PartnershipType,
		SponsorshipLevel:  req.SponsorshipLevel,
		SponsorshipAmount: req.SponsorshipAmount,
		Logo:              req.Logo,
	}); err != nil {
		return err
	}

	return nil
}

func (service *EventCompanyServiceImpl) Update(ctx context.Context, req model.EventCompanyRequest, ecID uint) error {
	eventCompany, err := service.Repository.FindOne(ecID)
	if err != nil {
		return err
	}

	if err = service.Repository.Update(ecID, model.EventCompany{
		BaseEntity:        builder.BuildBaseEntity(ctx, false, &eventCompany.BaseEntity),
		EventID:           eventCompany.EventID,
		Name:              req.Name,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		PartnershipType:   req.PartnershipType,
		SponsorshipLevel:  req.SponsorshipLevel,
		SponsorshipAmount: req.SponsorshipAmount,
		Logo:              req.Logo,
	}); err != nil {
		return err
	}

	return nil
}

func (service *EventCompanyServiceImpl) Delete(ecID uint) error {
	if _, err := service.Repository.FindOne(ecID); err != nil {
		return err
	}

	if err := service.Repository.Delete(ecID); err != nil {
		return err
	}

	return nil
}

func (service *EventCompanyServiceImpl) GetList(filter model.FilterEventCompany) (companies []model.EventCompany, err error) {
	if companies, err = service.Repository.FindAll(filter); err != nil {
		return
	}
	return
}

func (service *EventCompanyServiceImpl) GetDetail(ecID uint) (company model.EventCompany, err error) {
	if company, err = service.Repository.FindOne(ecID); err != nil {
		return
	}
	return
}
