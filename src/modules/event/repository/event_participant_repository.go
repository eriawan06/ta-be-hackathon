package repository

import (
	"be-sagara-hackathon/src/modules/event/model"
	e "be-sagara-hackathon/src/utils/errors"
	"gorm.io/gorm"
)

type EventParticipantRepository interface {
	Save(eventParticipant model.EventParticipant) error
	Update(eventParticipant model.EventParticipant, eventParticipantId uint) error
	Delete(eventParticipantId uint, deleteBy string) error
	FindAll() ([]model.EventParticipant, error)
	FindOne(eventParticipantID uint) (model.EventParticipant, error)
	FindOneByEventIDAndParticipantID(eventID, participantID uint) (model.EventParticipant, error)
}

type EventParticipantRepositoryImpl struct {
	DB *gorm.DB
}

func NewEventParticipantRepository(db *gorm.DB) EventParticipantRepository {
	return &EventParticipantRepositoryImpl{DB: db}
}

func (repository *EventParticipantRepositoryImpl) Save(eventParticipant model.EventParticipant) error {
	result := repository.DB.Create(&eventParticipant)
	return result.Error
}

func (repository *EventParticipantRepositoryImpl) Update(eventParticipant model.EventParticipant, eventParticipantId uint) error {
	result := repository.DB.Where("id = ?", eventParticipantId).Updates(&eventParticipant)
	return result.Error
}

func (repository *EventParticipantRepositoryImpl) Delete(eventParticipantId uint, deleteBy string) error {
	query := `UPDATE event_participants SET deleted_at=NOW(), deleted_by=? WHERE id=?`
	result := repository.DB.Exec(query, deleteBy, eventParticipantId)

	return result.Error
}

func (repository *EventParticipantRepositoryImpl) FindAll() ([]model.EventParticipant, error) {
	var eventParticipants []model.EventParticipant
	result := repository.DB.Where("deleted_at IS NULL").Find(&eventParticipants)
	return eventParticipants, result.Error
}

func (repository *EventParticipantRepositoryImpl) FindOne(eventParticipantID uint) (model.EventParticipant, error) {
	var eventParticipant model.EventParticipant
	result := repository.DB.Where("id = ?", eventParticipantID).First(&eventParticipant)
	if result.Error == gorm.ErrRecordNotFound {
		return eventParticipant, e.ErrDataNotFound
	}
	return eventParticipant, result.Error
}

func (repository *EventParticipantRepositoryImpl) FindOneByEventIDAndParticipantID(eventID, participantID uint) (model.EventParticipant, error) {
	var eventParticipant model.EventParticipant
	if err := repository.DB.Where("event_id=? AND participant_id=?", eventID, participantID).
		First(&eventParticipant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return eventParticipant, e.ErrDataNotFound
		}
		return eventParticipant, err
	}
	return eventParticipant, nil
}
