package service

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"context"
)

type EventAssessmentCriteriaService interface {
	Create(ctx context.Context, req model.EventAssessmentCriteriaRequest) error
	Update(ctx context.Context, req model.UpdateEventAssessmentCriteriaRequest, id uint) error
	Delete(id uint) error
	GetList(
		filter model.FilterEventAssessmentCriteria,
		pg *utils.PaginateQueryOffset,
	) (response model.ListEventAssessmentCriteriaResponse, err error)
	GetDetail(id uint) (criteria model.EventAssessmentCriteria, err error)
	GetActiveByEventID(eventID uint) (criteria []model.EventAssessmentCriteria, err error)
}

type EventAssessmentCriteriaServiceImpl struct {
	Repository      repository.EventAssessmentCriteriaRepository
	EventRepository repository.EventRepository
}

func NewEventAssessmentCriteriaService(
	repository repository.EventAssessmentCriteriaRepository,
	eventRepository repository.EventRepository,
) EventAssessmentCriteriaService {
	return &EventAssessmentCriteriaServiceImpl{Repository: repository, EventRepository: eventRepository}
}

func (service *EventAssessmentCriteriaServiceImpl) Create(ctx context.Context, req model.EventAssessmentCriteriaRequest) error {
	if _, err := service.EventRepository.FindOne(req.EventID); err != nil {
		return err
	}

	if err := service.Repository.Save(model.EventAssessmentCriteria{
		BaseEntity:    builder.BuildBaseEntity(ctx, true, nil),
		EventID:       req.EventID,
		Criteria:      req.Criteria,
		PercentageVal: req.PercentageVal,
		ScoreStart:    req.ScoreStart,
		ScoreEnd:      req.ScoreEnd,
		IsActive:      true,
	}); err != nil {
		return err
	}

	return nil
}

func (service *EventAssessmentCriteriaServiceImpl) Update(ctx context.Context, req model.UpdateEventAssessmentCriteriaRequest, id uint) error {
	criteria, err := service.Repository.FindOne(id)
	if err != nil {
		return err
	}

	if err = service.Repository.Update(id, model.EventAssessmentCriteria{
		BaseEntity:    builder.BuildBaseEntity(ctx, false, &criteria.BaseEntity),
		EventID:       criteria.EventID,
		Criteria:      req.Criteria,
		PercentageVal: req.PercentageVal,
		ScoreStart:    req.ScoreStart,
		ScoreEnd:      req.ScoreEnd,
		IsActive:      req.IsActive,
	}); err != nil {
		return err
	}

	return nil
}

func (service *EventAssessmentCriteriaServiceImpl) Delete(id uint) error {
	if _, err := service.Repository.FindOne(id); err != nil {
		return err
	}

	if err := service.Repository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (service *EventAssessmentCriteriaServiceImpl) GetList(
	filter model.FilterEventAssessmentCriteria,
	pg *utils.PaginateQueryOffset,
) (response model.ListEventAssessmentCriteriaResponse, err error) {
	response.EventCriteria, response.TotalItem, response.TotalPage, err = service.Repository.Find(filter, pg)
	if err != nil {
		return
	}
	return
}

func (service *EventAssessmentCriteriaServiceImpl) GetDetail(id uint) (criteria model.EventAssessmentCriteria, err error) {
	if criteria, err = service.Repository.FindOne(id); err != nil {
		return
	}
	return
}

func (service *EventAssessmentCriteriaServiceImpl) GetActiveByEventID(eventID uint) (criteria []model.EventAssessmentCriteria, err error) {
	if criteria, err = service.Repository.FindActiveByEventID(eventID); err != nil {
		return
	}
	return
}
