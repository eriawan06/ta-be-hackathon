package service

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/utils/common/builder"
	e "be-sagara-hackathon/src/utils/errors"
	"context"
)

type EventJudgeService interface {
	Create(ctx context.Context, request model.EventJudgeRequest) error
	Delete(ejID uint) error
	GetAll(filter model.FilterEventJudge) (judges []model.EventJudgeLite, err error)
	GetDetail(ejID uint) (judge model.EventJudge, err error)
}

type EventJudgeServiceImpl struct {
	Repository repository.EventJudgeRepository
}

func NewEventJudgeService(repository repository.EventJudgeRepository) EventJudgeService {
	return &EventJudgeServiceImpl{Repository: repository}
}

func (service *EventJudgeServiceImpl) Create(ctx context.Context, request model.EventJudgeRequest) error {
	existing, err := service.Repository.FindOneByJudgeIDAndEventID(request.JudgeID, request.EventID)
	if err != nil && err != e.ErrDataNotFound {
		return err
	}

	if existing.ID != 0 {
		return e.ErrEventJudgeExist
	}

	if err = service.Repository.Save(model.EventJudge{
		BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
		EventID:    request.EventID,
		JudgeID:    request.JudgeID,
	}); err != nil {
		return err
	}
	return nil
}

func (service *EventJudgeServiceImpl) Delete(ejID uint) error {
	if _, err := service.Repository.FindOne(ejID); err != nil {
		return err
	}
	if err := service.Repository.Delete(ejID); err != nil {
		return err
	}
	return nil
}

func (service *EventJudgeServiceImpl) GetAll(filter model.FilterEventJudge) (judges []model.EventJudgeLite, err error) {
	if judges, err = service.Repository.FindAll(filter); err != nil {
		return
	}
	return
}

func (service *EventJudgeServiceImpl) GetDetail(ejID uint) (judge model.EventJudge, err error) {
	if judge, err = service.Repository.FindOne(ejID); err != nil {
		return
	}
	return
}
