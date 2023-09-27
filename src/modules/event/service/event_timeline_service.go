package service

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/helper"
	"context"
)

type EventTimelineService interface {
	Create(ctx context.Context, request model.EventTimelineRequest) error
	Update(ctx context.Context, request model.EventTimelineRequest, etlID uint) error
	Delete(etlID uint) error
	GetList(filter model.FilterEventTimeline) ([]model.EventTimeline, error)
	GetDetail(etlID uint) (model.EventTimeline, error)
}

type EventTimelineServiceImpl struct {
	Repository      repository.EventTimelineRepository
	EventRepository repository.EventRepository
}

func NewEventTimelineService(
	repository repository.EventTimelineRepository,
	eventRepository repository.EventRepository,
) EventTimelineService {
	return &EventTimelineServiceImpl{Repository: repository, EventRepository: eventRepository}
}

func (service EventTimelineServiceImpl) Create(ctx context.Context, request model.EventTimelineRequest) error {
	if _, err := service.EventRepository.FindOne(request.EventID); err != nil {
		return err
	}

	startDate, err := helper.ParseDateStringToTime(request.StartDate)
	if err != nil {
		return err
	}
	endDate, err := helper.ParseDateStringToTime(request.EndDate)
	if err != nil {
		return err
	}
	if err = service.Repository.Save(model.EventTimeline{
		BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
		EventID:    request.EventID,
		Title:      request.Title,
		StartDate:  startDate,
		EndDate:    endDate,
		Note:       request.Note,
	}); err != nil {
		return err
	}

	return nil
}

func (service EventTimelineServiceImpl) Update(ctx context.Context, request model.EventTimelineRequest, etlID uint) error {
	//check event timeline
	existing, err := service.Repository.FindOne(etlID)
	if err != nil {
		return err
	}

	startDate, err := helper.ParseDateStringToTime(request.StartDate)
	if err != nil {
		return err
	}
	endDate, err := helper.ParseDateStringToTime(request.EndDate)
	if err != nil {
		return err
	}
	if err = service.Repository.Update(etlID, model.EventTimeline{
		BaseEntity: builder.BuildBaseEntity(ctx, false, &existing.BaseEntity),
		EventID:    existing.EventID,
		Title:      request.Title,
		StartDate:  startDate,
		EndDate:    endDate,
		Note:       request.Note,
	}); err != nil {
		return err
	}

	return nil
}

func (service EventTimelineServiceImpl) Delete(etlID uint) error {
	//check event timeline
	_, err := service.Repository.FindOne(etlID)
	if err != nil {
		return err
	}

	if err = service.Repository.Delete(etlID); err != nil {
		return err
	}
	return nil
}

func (service EventTimelineServiceImpl) GetList(filter model.FilterEventTimeline) ([]model.EventTimeline, error) {
	timelines, err := service.Repository.FindAll(filter)
	if err != nil {
		return timelines, err
	}
	return timelines, nil
}

func (service EventTimelineServiceImpl) GetDetail(etlID uint) (model.EventTimeline, error) {
	var timeline model.EventTimeline
	timeline, err := service.Repository.FindOne(etlID)
	if err != nil {
		return timeline, err
	}
	return timeline, nil
}
