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

type EventFaqRepository interface {
	Save(req model.EventFaq) (err error)
	Update(id uint, req model.EventFaq) (err error)
	Delete(id uint) (err error)
	Find(
		filter model.FilterEventFaq,
		pg *utils.PaginateQueryOffset,
	) (faqs []model.EventFaq, totalData, totalPage int64, err error)
	FindOne(id uint) (faq model.EventFaq, err error)
	FindManyByEventID(eventID uint) (faqs []model.EventFaq, err error)
}

type EventFaqRepositoryImpl struct {
	DB *gorm.DB
}

func NewEventFaqRepository(db *gorm.DB) EventFaqRepository {
	return &EventFaqRepositoryImpl{DB: db}
}

func (repository *EventFaqRepositoryImpl) Save(req model.EventFaq) (err error) {
	if err = repository.DB.Create(&req).Error; err != nil {
		return
	}
	return
}

func (repository *EventFaqRepositoryImpl) Update(id uint, req model.EventFaq) (err error) {
	if err = repository.DB.Select("*").Where("id = ?", id).Updates(&req).Error; err != nil {
		return
	}
	return
}

func (repository *EventFaqRepositoryImpl) Delete(id uint) (err error) {
	if err = repository.DB.Delete(&model.EventFaq{}, id).Error; err != nil {
		return
	}
	return
}

func (repository *EventFaqRepositoryImpl) Find(
	filter model.FilterEventFaq,
	pg *utils.PaginateQueryOffset,
) (faqs []model.EventFaq, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterEventFaq(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&faqs).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalEventFaq(&filter)
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

func (repository *EventFaqRepositoryImpl) getTotalEventFaq(filter *model.FilterEventFaq) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterEventFaq(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.EventFaq{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterEventFaq(filter model.FilterEventFaq) (where []string, whereVal []interface{}) {
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

func (repository *EventFaqRepositoryImpl) FindOne(id uint) (faq model.EventFaq, err error) {
	if err = repository.DB.Where("id=?", id).First(&faq).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *EventFaqRepositoryImpl) FindManyByEventID(eventID uint) (faqs []model.EventFaq, err error) {
	if err = repository.DB.Where("event_id=? AND is_active=true", eventID).Find(&faqs).Error; err != nil {
		return
	}
	return
}
