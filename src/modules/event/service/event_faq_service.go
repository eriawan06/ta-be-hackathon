package service

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"context"
)

type EventFaqService interface {
	Create(ctx context.Context, req model.EventFaqRequest) error
	Update(ctx context.Context, req model.UpdateEventFaqRequest, id uint) error
	Delete(id uint) error
	GetList(
		filter model.FilterEventFaq,
		pg *utils.PaginateQueryOffset,
	) (response model.ListEventFaqResponse, err error)
	GetDetail(id uint) (faq model.EventFaq, err error)
}

type EventFaqServiceImpl struct {
	Repository      repository.EventFaqRepository
	EventRepository repository.EventRepository
}

func NewEventFaqService(
	repository repository.EventFaqRepository,
	eventRepository repository.EventRepository,
) EventFaqService {
	return &EventFaqServiceImpl{Repository: repository, EventRepository: eventRepository}
}

func (service *EventFaqServiceImpl) Create(ctx context.Context, req model.EventFaqRequest) error {
	if _, err := service.EventRepository.FindOne(req.EventID); err != nil {
		return err
	}

	if err := service.Repository.Save(model.EventFaq{
		BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
		EventID:    req.EventID,
		Title:      req.Title,
		Note:       req.Note,
		IsActive:   true,
	}); err != nil {
		return err
	}

	return nil
}

func (service *EventFaqServiceImpl) Update(ctx context.Context, req model.UpdateEventFaqRequest, id uint) error {
	eventFaq, err := service.Repository.FindOne(id)
	if err != nil {
		return err
	}

	if err = service.Repository.Update(id, model.EventFaq{
		BaseEntity: builder.BuildBaseEntity(ctx, false, &eventFaq.BaseEntity),
		EventID:    eventFaq.EventID,
		Title:      req.Title,
		Note:       req.Note,
		IsActive:   req.IsActive,
	}); err != nil {
		return err
	}

	return nil
}

func (service *EventFaqServiceImpl) Delete(id uint) error {
	if _, err := service.Repository.FindOne(id); err != nil {
		return err
	}

	if err := service.Repository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (service *EventFaqServiceImpl) GetList(
	filter model.FilterEventFaq,
	pg *utils.PaginateQueryOffset,
) (response model.ListEventFaqResponse, err error) {
	faqs, totalData, totalPage, err := service.Repository.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.EventFaqLite
	for _, v := range faqs {
		responseData = append(responseData, model.EventFaqLite{
			ID:        v.ID,
			Title:     v.Title,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			IsActive:  v.IsActive,
		})
	}
	response.EventFaqs = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *EventFaqServiceImpl) GetDetail(id uint) (faq model.EventFaq, err error) {
	if faq, err = service.Repository.FindOne(id); err != nil {
		return
	}
	return
}
