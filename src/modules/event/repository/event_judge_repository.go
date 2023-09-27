package repository

import (
	"be-sagara-hackathon/src/modules/event/model"
	e "be-sagara-hackathon/src/utils/errors"
	"gorm.io/gorm"
)

type EventJudgeRepository interface {
	Save(ej model.EventJudge) error
	Delete(ejID uint) error
	FindAll(filter model.FilterEventJudge) (judges []model.EventJudgeLite, err error)
	FindOne(ejID uint) (judge model.EventJudge, err error)
	FindManyByEventID(eventID uint) ([]model.EventJudgeLite, error)
	FindOneByJudgeIDAndEventID(judgeID, eventID uint) (judge model.EventJudge, err error)
}

type EventJudgeRepositoryImpl struct {
	DB *gorm.DB
}

func NewEventJudgeRepository(db *gorm.DB) EventJudgeRepository {
	return &EventJudgeRepositoryImpl{DB: db}
}

func (repository *EventJudgeRepositoryImpl) Save(ej model.EventJudge) error {
	if err := repository.DB.Create(&ej).Error; err != nil {
		return err
	}
	return nil
}

func (repository *EventJudgeRepositoryImpl) Delete(ejID uint) error {
	if err := repository.DB.Delete(&model.EventJudge{}, ejID).Error; err != nil {
		return err
	}
	return nil
}

func (repository *EventJudgeRepositoryImpl) FindAll(filter model.FilterEventJudge) (judges []model.EventJudgeLite, err error) {
	query := `
		SELECT ej.*, u.name as judge_name, oc.name as judge_occupation, 
		       u.institution as judge_institution, u.avatar as judge_avatar
		FROM event_judges ej
		INNER JOIN events e on e.id = ej.event_id
		INNER JOIN users u on u.id = ej.judge_id
		INNER JOIN occupations oc on oc.id = u.occupation_id
	`

	if filter.EventID > 0 {
		query += ` WHERE ej.event_id = ?`
	}

	if err = repository.DB.Raw(query, filter.EventID).Scan(&judges).Error; err != nil {
		return
	}
	return
}

func (repository *EventJudgeRepositoryImpl) FindOne(ejID uint) (judge model.EventJudge, err error) {
	if err = repository.DB.Where("id=?", ejID).
		Preload("Judge").
		Preload("Judge.Occupation").
		First(&judge).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *EventJudgeRepositoryImpl) FindManyByEventID(eventID uint) ([]model.EventJudgeLite, error) {
	judges := []model.EventJudgeLite{}
	query := `
		SELECT ej.*, u.name as judge_name, oc.name as judge_occupation, 
		       u.institution as judge_institution, u.avatar as judge_avatar
		FROM event_judges ej
		INNER JOIN events e on e.id = ej.event_id
		INNER JOIN users u on u.id = ej.judge_id
		INNER JOIN occupations oc on oc.id = u.occupation_id
		WHERE ej.event_id=?
	`
	if err := repository.DB.Raw(query, eventID).Scan(&judges).Error; err != nil {
		return judges, err
	}
	return judges, nil
}

func (repository *EventJudgeRepositoryImpl) FindOneByJudgeIDAndEventID(judgeID, eventID uint) (judge model.EventJudge, err error) {
	if err = repository.DB.Where("judge_id=? AND event_id=?", judgeID, eventID).
		First(&judge).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}
