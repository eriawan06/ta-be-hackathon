package repository

import (
	"be-sagara-hackathon/src/modules/event/model"
	e "be-sagara-hackathon/src/utils/errors"
	"gorm.io/gorm"
)

type EventTimelineRepository interface {
	Save(etl model.EventTimeline) error
	Update(etlID uint, etl model.EventTimeline) error
	Delete(etlID uint) error
	FindAll(filter model.FilterEventTimeline) ([]model.EventTimeline, error)
	FindOne(etlID uint) (model.EventTimeline, error)
	FindManyByEventID(eventID uint) ([]model.EventTimeline, error)
}

type EventTimelineRepositoryImpl struct {
	DB *gorm.DB
}

func NewEventTimelineRepository(db *gorm.DB) EventTimelineRepository {
	return &EventTimelineRepositoryImpl{DB: db}
}

func (repository EventTimelineRepositoryImpl) Save(etl model.EventTimeline) error {
	if err := repository.DB.Create(&etl).Error; err != nil {
		return err
	}
	return nil
}

func (repository EventTimelineRepositoryImpl) Update(etlID uint, etl model.EventTimeline) error {
	if err := repository.DB.Select("*").Where("id=?", etlID).Updates(&etl).Error; err != nil {
		return err
	}
	return nil
}

func (repository EventTimelineRepositoryImpl) Delete(etlID uint) error {
	if err := repository.DB.Delete(&model.EventTimeline{}, etlID).Error; err != nil {
		return err
	}
	return nil
}

func (repository EventTimelineRepositoryImpl) FindAll(filter model.FilterEventTimeline) ([]model.EventTimeline, error) {
	var eventTimelines []model.EventTimeline
	if err := repository.DB.Where("event_id=?", filter.EventID).
		Find(&eventTimelines).Error; err != nil {
		return eventTimelines, err
	}
	return eventTimelines, nil
}

func (repository EventTimelineRepositoryImpl) FindOne(etlID uint) (model.EventTimeline, error) {
	var eventTimeline model.EventTimeline
	if err := repository.DB.Where("id = ?", etlID).First(&eventTimeline).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return eventTimeline, err
	}

	return eventTimeline, nil
}

func (repository EventTimelineRepositoryImpl) FindManyByEventID(eventID uint) ([]model.EventTimeline, error) {
	var eventTimelines []model.EventTimeline
	if err := repository.DB.Where("event_id=?", eventID).Find(&eventTimelines).Error; err != nil {
		return eventTimelines, err
	}
	return eventTimelines, nil
}
