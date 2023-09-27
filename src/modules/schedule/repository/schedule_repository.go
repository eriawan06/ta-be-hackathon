package repository

import (
	"be-sagara-hackathon/src/modules/schedule/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type ScheduleRepository interface {
	Save(req model.Schedule) (schedule model.Schedule, err error)
	SaveScheduleTeam(req model.ScheduleTeam) (err error)
	Update(id uint, req model.Schedule) (err error)
	Delete(id uint) (err error)
	DeleteScheduleTeam(id, teamID uint) (err error)
	Find(
		filter model.FilterSchedule,
		pg *utils.PaginateQueryOffset,
	) (schedule []model.ScheduleLite, totalData, totalPage int64, err error)
	FindDetail(id uint) (schedule model.ScheduleDetail, err error)
	FindOne(id uint) (schedule model.Schedule, err error)
	FindByEventIDAndTeamID(eventID, teamID uint) (schedules []model.ScheduleLite2, err error)
}

type ScheduleRepositoryImpl struct {
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &ScheduleRepositoryImpl{DB: db}
}

func (repository *ScheduleRepositoryImpl) Save(req model.Schedule) (schedule model.Schedule, err error) {
	if err = repository.DB.Omit("Event").Omit("Mentor").Create(&req).Error; err != nil {
		return
	}
	schedule = req
	return
}

func (repository *ScheduleRepositoryImpl) SaveScheduleTeam(req model.ScheduleTeam) (err error) {
	if err = repository.DB.Create(&req).Error; err != nil {
		return
	}
	return
}

func (repository *ScheduleRepositoryImpl) Update(id uint, req model.Schedule) (err error) {
	if err = repository.DB.Select("*").Where("id = ?", id).Updates(&req).Error; err != nil {
		return
	}
	return
}

func (repository *ScheduleRepositoryImpl) Delete(id uint) (err error) {
	tx := repository.DB.Begin()
	if err = tx.Delete(&model.Schedule{}, "id=?", id).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Delete(&model.ScheduleTeam{}, "schedule_id=?", id).Error; err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

func (repository *ScheduleRepositoryImpl) DeleteScheduleTeam(id, teamID uint) (err error) {
	tx := repository.DB.Begin()
	if err = tx.Delete(&model.ScheduleTeam{}, "schedule_id=? AND team_id=?", id, teamID).Error; err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

func (repository *ScheduleRepositoryImpl) Find(
	filter model.FilterSchedule,
	pg *utils.PaginateQueryOffset,
) (schedule []model.ScheduleLite, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilter(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("schedules as s").
		Select("s.id, s.event_id, s.title, s.held_on, u.name as mentor_name").
		Joins("inner join users u on u.id = s.mentor_id").
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&schedule).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalSchedule(&filter)
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

func (repository *ScheduleRepositoryImpl) getTotalSchedule(filter *model.FilterSchedule) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilter(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("schedules as s").
		Joins("inner join users u on u.id = s.mentor_id").
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilter(filter model.FilterSchedule) (where []string, whereVal []interface{}) {
	if filter.Search != "" {
		filter.Search = strings.ToLower(filter.Search)
		where = append(where, "(LOWER(s.title) LIKE @q OR LOWER(u.name) LIKE @q)")
		whereVal = append(whereVal, sql.Named("q", "%"+filter.Search+"%"))
	}

	if filter.EventID != 0 {
		where = append(where, "s.event_id = @event")
		whereVal = append(whereVal, sql.Named("event", filter.EventID))
	}

	if filter.HeldOn != "" {
		where = append(where, "s.held_on = @held")
		whereVal = append(whereVal, sql.Named("held", filter.HeldOn))
	}

	if filter.MentorID != 0 {
		where = append(where, "s.mentor_id = @mentor")
		whereVal = append(whereVal, sql.Named("mentor", filter.MentorID))
	}

	return
}

func (repository *ScheduleRepositoryImpl) FindDetail(id uint) (schedule model.ScheduleDetail, err error) {
	if err = repository.DB.Table("schedules as s").
		Select(`s.id, s.event_id, s.title, s.held_on, 
			u.name as mentor_name, o.name as mentor_occupation, u.institution as mentor_institution,
			u.avatar as mentor_avatar`).
		Joins("inner join users u on u.id = s.mentor_id").
		Joins("inner join occupations o on o.id = u.occupation_id").
		Where("s.id=?", id).
		First(&schedule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *ScheduleRepositoryImpl) FindOne(id uint) (schedule model.Schedule, err error) {
	if err = repository.DB.Where("id=?", id).First(&schedule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *ScheduleRepositoryImpl) FindByEventIDAndTeamID(eventID, teamID uint) (schedules []model.ScheduleLite2, err error) {
	if err = repository.DB.Table("schedules as s").
		Select(`s.id, s.title, s.held_on`).
		Joins("inner join schedule_teams st on st.schedule_id = s.id").
		Where("s.event_id=? AND st.team_id=?", eventID, teamID).
		Find(&schedules).Error; err != nil {
		return
	}
	return
}
