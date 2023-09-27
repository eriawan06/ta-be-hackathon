package repository

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type EventRepository interface {
	Save(event *model.Event) error
	Update(event model.Event, eventID uint) error
	Delete(eventID uint, deleteBy string) error
	Find(
		filter model.FilterEvent,
		pg *utils.PaginateQueryOffset,
	) (events []model.Event, totalData, totalPage int64, err error)
	FindOne(eventID uint) (event model.Event, err error)
	FindLatest() (event model.Event, err error)
}

type EventRepositoryImpl struct {
	DB *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &EventRepositoryImpl{DB: db}
}

func (repository *EventRepositoryImpl) Save(event *model.Event) error {
	if err := repository.DB.Create(&event).Error; err != nil {
		return err
	}
	return nil
}

func (repository *EventRepositoryImpl) Update(event model.Event, eventID uint) error {
	if err := repository.DB.Select("*").
		Where("id = ?", eventID).
		Updates(&event).Error; err != nil {
		return err
	}
	return nil
}

func (repository *EventRepositoryImpl) Delete(eventID uint, deleteBy string) error {
	query := `UPDATE events SET deleted_at=NOW(), deleted_by=? WHERE id=?`
	if err := repository.DB.Exec(query, deleteBy, eventID).Error; err != nil {
		return err
	}

	return nil
}

func (repository *EventRepositoryImpl) Find(
	filter model.FilterEvent,
	pg *utils.PaginateQueryOffset,
) (events []model.Event, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilter(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&events).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalEvent(&filter)
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

func (repository *EventRepositoryImpl) getTotalEvent(filter *model.FilterEvent) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilter(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.Event{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilter(filter model.FilterEvent) (where []string, whereVal []interface{}) {
	where = append(where, "deleted_at IS NULL")

	if filter.Search != "" {
		where = append(where, "(LOWER(name) LIKE @keyword OR reg_fee = @keyword)")
		whereVal = append(whereVal, sql.Named("keyword", "%"+filter.Search+"%"))
	}

	if filter.Status != "" {
		where = append(where, "status = @status")
		whereVal = append(whereVal, sql.Named("status", filter.Status))
	}

	if filter.StartDate != "" {
		where = append(where, "date(start_date) >= @start")
		whereVal = append(whereVal, sql.Named("start", filter.StartDate))
	}

	if filter.EndDate != "" {
		where = append(where, "date(end_date) <= @end")
		whereVal = append(whereVal, sql.Named("end", filter.EndDate))
	}

	return
}

func (repository *EventRepositoryImpl) FindOne(eventID uint) (event model.Event, err error) {
	if err = repository.DB.Where("id = ?", eventID).
		Preload("Mentors").Preload("Mentors.Mentor").
		Preload("Judges").Preload("Judges.Judge").
		Preload("Companies").
		Preload("Timelines").
		Preload("Rules").
		Preload("FAQs").
		First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *EventRepositoryImpl) FindLatest() (event model.Event, err error) {
	if err = repository.DB.Where("deleted_at IS NULL").
		Order("id desc").
		First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}
