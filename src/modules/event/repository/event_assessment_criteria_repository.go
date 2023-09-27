package repository

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type EventAssessmentCriteriaRepository interface {
	Save(req model.EventAssessmentCriteria) (err error)
	Update(id uint, req model.EventAssessmentCriteria) (err error)
	Delete(id uint) (err error)
	Find(
		filter model.FilterEventAssessmentCriteria,
		pg *utils.PaginateQueryOffset,
	) (criteria []model.EventAssessmentCriteria, totalData, totalPage int64, err error)
	FindOne(id uint) (criteria model.EventAssessmentCriteria, err error)
	FindActiveByEventID(eventID uint) (criteria []model.EventAssessmentCriteria, err error)
}

type EventAssessmentCriteriaRepositoryImpl struct {
	DB *gorm.DB
}

func NewEventAssessmentCriteriaRepository(db *gorm.DB) EventAssessmentCriteriaRepository {
	return &EventAssessmentCriteriaRepositoryImpl{DB: db}
}

func (repository *EventAssessmentCriteriaRepositoryImpl) Save(req model.EventAssessmentCriteria) (err error) {
	if err = repository.DB.Create(&req).Error; err != nil {
		return
	}
	return
}

func (repository *EventAssessmentCriteriaRepositoryImpl) Update(id uint, req model.EventAssessmentCriteria) (err error) {
	if err = repository.DB.Select("*").Where("id = ?", id).Updates(&req).Error; err != nil {
		return
	}
	return
}

func (repository *EventAssessmentCriteriaRepositoryImpl) Delete(id uint) (err error) {
	if err = repository.DB.Delete(&model.EventAssessmentCriteria{}, id).Error; err != nil {
		return
	}
	return
}

func (repository *EventAssessmentCriteriaRepositoryImpl) Find(
	filter model.FilterEventAssessmentCriteria,
	pg *utils.PaginateQueryOffset,
) (criteria []model.EventAssessmentCriteria, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterEventAssessmentCriteria(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&criteria).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalEventAssessmentCriteria(&filter)
	if err != nil {
		return
	}

	if pg.Limit > 0 {
		totalPage = int64(math.Ceil(float64(totalData) / float64(pg.Limit)))
	} else {
		totalPage = 1
	}

	return
}

func (repository *EventAssessmentCriteriaRepositoryImpl) getTotalEventAssessmentCriteria(filter *model.FilterEventAssessmentCriteria) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterEventAssessmentCriteria(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.EventAssessmentCriteria{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterEventAssessmentCriteria(filter model.FilterEventAssessmentCriteria) (where []string, whereVal []interface{}) {
	if filter.EventID > 0 {
		where = append(where, "event_id = ?")
		whereVal = append(whereVal, filter.EventID)
	}

	if filter.Status != "" {
		isActive := false
		if filter.Status == "active" {
			isActive = true
		}
		where = append(where, "is_active = ?")
		whereVal = append(whereVal, isActive)
	}

	return
}

func (repository *EventAssessmentCriteriaRepositoryImpl) FindOne(id uint) (criteria model.EventAssessmentCriteria, err error) {
	if err = repository.DB.Where("id=?", id).First(&criteria).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *EventAssessmentCriteriaRepositoryImpl) FindActiveByEventID(eventID uint) (criteria []model.EventAssessmentCriteria, err error) {
	if err = repository.DB.Where("event_id=? AND is_active=true AND deleted_at is null", eventID).
		Order("id asc").
		Find(&criteria).Error; err != nil {
		return
	}
	return
}
