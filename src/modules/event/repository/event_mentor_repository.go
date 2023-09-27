package repository

import (
	"be-sagara-hackathon/src/modules/event/model"
	e "be-sagara-hackathon/src/utils/errors"
	"gorm.io/gorm"
)

type EventMentorRepository interface {
	Save(em model.EventMentor) error
	Delete(emID uint) error
	FindAll(filter model.FilterEventMentor) (mentors []model.EventMentorLite, err error)
	FindOne(emID uint) (mentor model.EventMentor, err error)
	FindManyByEventID(eventID uint) (mentors []model.EventMentorLite, err error)
	FindOneByMentorIDAndEventID(mentorID, eventID uint) (mentor model.EventMentor, err error)
}

type EventMentorRepositoryImpl struct {
	DB *gorm.DB
}

func NewEventMentorRepository(db *gorm.DB) EventMentorRepository {
	return &EventMentorRepositoryImpl{DB: db}
}

func (repository *EventMentorRepositoryImpl) Save(em model.EventMentor) error {
	if err := repository.DB.Create(&em).Error; err != nil {
		return err
	}
	return nil
}

func (repository *EventMentorRepositoryImpl) Delete(emID uint) error {
	if err := repository.DB.Delete(&model.EventMentor{}, emID).Error; err != nil {
		return err
	}
	return nil
}

func (repository *EventMentorRepositoryImpl) FindAll(filter model.FilterEventMentor) (mentors []model.EventMentorLite, err error) {
	query := `
		SELECT em.*, u.name as mentor_name, oc.name as mentor_occupation, 
		       u.institution as mentor_institution, u.avatar as mentor_avatar
		FROM event_mentors em
		INNER JOIN events e on e.id = em.event_id
		INNER JOIN users u on u.id = em.mentor_id
		INNER JOIN occupations oc on oc.id = u.occupation_id
	`

	if filter.EventID > 0 {
		query += ` WHERE em.event_id = ?`
	}

	if err = repository.DB.Raw(query, filter.EventID).Scan(&mentors).Error; err != nil {
		return
	}
	return
}

func (repository *EventMentorRepositoryImpl) FindOne(emID uint) (mentor model.EventMentor, err error) {
	if err = repository.DB.Where("id=?", emID).
		Preload("Mentor").
		Preload("Mentor.Occupation").
		First(&mentor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *EventMentorRepositoryImpl) FindManyByEventID(eventID uint) ([]model.EventMentorLite, error) {
	mentors := []model.EventMentorLite{}
	query := `
		SELECT em.*, u.name as mentor_name, oc.name as mentor_occupation, 
		       u.institution as mentor_institution, u.avatar as mentor_avatar
		FROM event_mentors em
		INNER JOIN events e on e.id = em.event_id
		INNER JOIN users u on u.id = em.mentor_id
		INNER JOIN occupations oc on oc.id = u.occupation_id
		WHERE em.event_id=?
	`
	if err := repository.DB.Raw(query, eventID).Scan(&mentors).Error; err != nil {
		return mentors, err
	}
	return mentors, nil
}

func (repository *EventMentorRepositoryImpl) FindOneByMentorIDAndEventID(mentorID, eventID uint) (mentor model.EventMentor, err error) {
	if err = repository.DB.Where("mentor_id=? AND event_id=?", mentorID, eventID).
		First(&mentor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}
