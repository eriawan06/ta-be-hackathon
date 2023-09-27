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

type EventRuleRepository interface {
	Save(req model.EventRule) (err error)
	Update(id uint, req model.EventRule) (err error)
	Delete(id uint) (err error)
	Find(
		filter model.FilterEventRule,
		pg *utils.PaginateQueryOffset,
	) (rules []model.EventRule, totalData, totalPage int64, err error)
	FindOne(id uint) (rule model.EventRule, err error)
	FindActiveByEventID(eventID uint) (rules []model.EventRule, err error)
}

type EventRuleRepositoryImpl struct {
	DB *gorm.DB
}

func NewEventRuleRepository(db *gorm.DB) EventRuleRepository {
	return &EventRuleRepositoryImpl{DB: db}
}

func (repository *EventRuleRepositoryImpl) Save(req model.EventRule) (err error) {
	if err = repository.DB.Create(&req).Error; err != nil {
		return
	}
	return
}

func (repository *EventRuleRepositoryImpl) Update(id uint, req model.EventRule) (err error) {
	if err = repository.DB.Select("*").Where("id = ?", id).Updates(&req).Error; err != nil {
		return
	}
	return
}

func (repository *EventRuleRepositoryImpl) Delete(id uint) (err error) {
	if err = repository.DB.Delete(&model.EventRule{}, id).Error; err != nil {
		return
	}
	return
}

func (repository *EventRuleRepositoryImpl) Find(
	filter model.FilterEventRule,
	pg *utils.PaginateQueryOffset,
) (rules []model.EventRule, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterEventRule(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&rules).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalEventRule(&filter)
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

func (repository *EventRuleRepositoryImpl) getTotalEventRule(filter *model.FilterEventRule) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterEventRule(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.EventRule{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterEventRule(filter model.FilterEventRule) (where []string, whereVal []interface{}) {
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

func (repository *EventRuleRepositoryImpl) FindOne(id uint) (rule model.EventRule, err error) {
	if err = repository.DB.Where("id=?", id).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *EventRuleRepositoryImpl) FindActiveByEventID(eventID uint) (rules []model.EventRule, err error) {
	if err = repository.DB.Where("event_id=? AND is_active=true AND deleted_at is null", eventID).
		Order("id asc").
		Find(&rules).Error; err != nil {
		return
	}
	return
}
