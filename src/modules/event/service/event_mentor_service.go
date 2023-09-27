package service

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/utils/common/builder"
	e "be-sagara-hackathon/src/utils/errors"
	"context"
)

type EventMentorService interface {
	Create(ctx context.Context, request model.EventMentorRequest) error
	Delete(emID uint) error
	GetAll(filter model.FilterEventMentor) (mentors []model.EventMentorLite, err error)
	GetDetail(emID uint) (mentor model.EventMentor, err error)
}

type EventMentorServiceImpl struct {
	Repository repository.EventMentorRepository
}

func NewEventMentorService(repository repository.EventMentorRepository) EventMentorService {
	return &EventMentorServiceImpl{Repository: repository}
}

func (service *EventMentorServiceImpl) Create(ctx context.Context, request model.EventMentorRequest) error {
	existing, err := service.Repository.FindOneByMentorIDAndEventID(request.MentorID, request.EventID)
	if err != nil && err != e.ErrDataNotFound {
		return err
	}

	if existing.ID != 0 {
		return e.ErrEventMentorExist
	}

	if err = service.Repository.Save(model.EventMentor{
		BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
		EventID:    request.EventID,
		MentorID:   request.MentorID,
	}); err != nil {
		return err
	}
	return nil
}

func (service *EventMentorServiceImpl) Delete(emID uint) error {
	if _, err := service.Repository.FindOne(emID); err != nil {
		return err
	}
	if err := service.Repository.Delete(emID); err != nil {
		return err
	}
	return nil
}

func (service *EventMentorServiceImpl) GetAll(filter model.FilterEventMentor) (mentors []model.EventMentorLite, err error) {
	if mentors, err = service.Repository.FindAll(filter); err != nil {
		return
	}
	return
}

func (service *EventMentorServiceImpl) GetDetail(emID uint) (mentor model.EventMentor, err error) {
	if mentor, err = service.Repository.FindOne(emID); err != nil {
		return
	}
	return
}
