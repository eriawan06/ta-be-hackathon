package service

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"context"
)

type EventRuleService interface {
	Create(ctx context.Context, req model.EventRuleRequest) error
	Update(ctx context.Context, req model.UpdateEventRuleRequest, erID uint) error
	Delete(erID uint) error
	GetList(
		filter model.FilterEventRule,
		pg *utils.PaginateQueryOffset,
	) (response model.ListEventRuleResponse, err error)
	GetDetail(erID uint) (rule model.EventRule, err error)
	GetActiveByEventID(eventID uint) (rules []model.EventRule, err error)
}

type EventRuleServiceImpl struct {
	Repository      repository.EventRuleRepository
	EventRepository repository.EventRepository
}

func NewEventRuleService(
	repository repository.EventRuleRepository,
	eventRepository repository.EventRepository,
) EventRuleService {
	return &EventRuleServiceImpl{Repository: repository, EventRepository: eventRepository}
}

func (service *EventRuleServiceImpl) Create(ctx context.Context, req model.EventRuleRequest) error {
	if _, err := service.EventRepository.FindOne(req.EventID); err != nil {
		return err
	}

	if err := service.Repository.Save(model.EventRule{
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

func (service *EventRuleServiceImpl) Update(ctx context.Context, req model.UpdateEventRuleRequest, erID uint) error {
	eventRule, err := service.Repository.FindOne(erID)
	if err != nil {
		return err
	}

	if err = service.Repository.Update(erID, model.EventRule{
		BaseEntity: builder.BuildBaseEntity(ctx, false, &eventRule.BaseEntity),
		EventID:    eventRule.EventID,
		Title:      req.Title,
		Note:       req.Note,
		IsActive:   req.IsActive,
	}); err != nil {
		return err
	}

	return nil
}

func (service *EventRuleServiceImpl) Delete(erID uint) error {
	if _, err := service.Repository.FindOne(erID); err != nil {
		return err
	}

	if err := service.Repository.Delete(erID); err != nil {
		return err
	}

	return nil
}

func (service *EventRuleServiceImpl) GetList(
	filter model.FilterEventRule,
	pg *utils.PaginateQueryOffset,
) (response model.ListEventRuleResponse, err error) {
	rules, totalData, totalPage, err := service.Repository.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.EventRuleLite
	for _, v := range rules {
		responseData = append(responseData, model.EventRuleLite{
			ID:        v.ID,
			Title:     v.Title,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			IsActive:  v.IsActive,
		})
	}
	response.EventRules = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *EventRuleServiceImpl) GetDetail(erID uint) (rule model.EventRule, err error) {
	if rule, err = service.Repository.FindOne(erID); err != nil {
		return
	}
	return
}

func (service *EventRuleServiceImpl) GetActiveByEventID(eventID uint) (rules []model.EventRule, err error) {
	if rules, err = service.Repository.FindActiveByEventID(eventID); err != nil {
		return
	}
	return
}
