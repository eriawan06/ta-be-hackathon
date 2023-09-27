package repository

import (
	"be-sagara-hackathon/src/modules/team/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

type TeamRepository interface {
	Save(eventID uint, team model.Team, member model.TeamMember) (newTeam model.Team, err error)
	Update(id uint, team model.Team) error
	Delete(id uint, deletedBy string) error
	FindOne(id uint) (team model.Team, err error)
	FindAll(
		filter model.FilterTeam,
		pg *utils.PaginateQueryOffset,
	) (teams []model.TeamLite, totalData, totalPage int64, err error)
	FindDetail(id uint) (team model.TeamDetail, err error)
	FindDetail2(id, eventID, participantID uint, includeMembers bool) (team model.TeamDetail2, err error)
	FindByIDAndEventID(id, eventID uint) (team model.TeamDetail2, err error)
}

type TeamRepositoryImpl struct {
	DB *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &TeamRepositoryImpl{DB: db}
}

func (repository *TeamRepositoryImpl) Save(eventID uint, team model.Team, member model.TeamMember) (newTeam model.Team, err error) {
	tx := repository.DB.Begin()
	if err = tx.Model(&model.Team{}).Omit("Participant").Create(&team).Error; err != nil {

		var mySqlErr *mysql.MySQLError
		if errors.As(err, &mySqlErr) && mySqlErr.Number == 1062 {
			if strings.Contains(mySqlErr.Message, "idx_unique_team_code") {
				err = e.ErrTeamCodeAlreadyExists
			} else if strings.Contains(mySqlErr.Message, "idx_unique_team_name") {
				err = e.ErrTeamNameAlreadyExists
			}
		}

		tx.Rollback()
		return
	}

	if err = tx.Model(&model.TeamEvent{}).Create(&model.TeamEvent{
		TeamID:  team.ID,
		EventID: eventID,
	}).Error; err != nil {
		tx.Rollback()
		return
	}

	member.TeamID = team.ID
	if err = tx.Model(&model.TeamMember{}).
		Omit("Team").Omit("Participant").
		Create(&member).Error; err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	newTeam = team
	return
}

func (repository *TeamRepositoryImpl) Update(id uint, team model.Team) error {
	if err := repository.DB.Select("*").
		Where("id=?", id).
		Updates(&team).Error; err != nil {
		return err
	}
	return nil
}

func (repository *TeamRepositoryImpl) Delete(id uint, deletedBy string) error {
	if err := repository.DB.Model(&model.Team{}).Where("id=?", id).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": deletedBy,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (repository *TeamRepositoryImpl) FindOne(id uint) (team model.Team, err error) {
	if err = repository.DB.Where("id=?", id).First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *TeamRepositoryImpl) FindAll(
	filter model.FilterTeam,
	pg *utils.PaginateQueryOffset,
) (teams []model.TeamLite, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterTeam(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	qSelect := `t.id, t.code, t.name, te.event_id, t.is_active, t.avatar, count(tm.id) as num_of_member, u.name as participant_name`
	if filter.EventID != 0 && filter.TeamRequestParticipantID != 0 {
		qSelect = fmt.Sprintf(`%s,
			IF((select count(tr.id)
			from team_requests tr
			where tr.team_id = t.id
				and tr.participant_id = %d
				and tr.event_id = %d
				and tr.status='sent') = 0, false, true) as is_requested`, qSelect, filter.TeamRequestParticipantID, filter.EventID)
	}

	db := repository.DB.Table("teams as t").
		Select(qSelect).
		Joins("inner join team_events te on te.team_id = t.id").
		Joins("inner join team_members tm on tm.team_id = t.id").
		Joins("inner join participants p on p.id = t.participant_id").
		Joins("inner join users u on u.id = p.user_id")

	if filter.ScheduleID != 0 {
		db.Joins("left join schedule_teams st on st.team_id = t.id")
	}

	if err = db.Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Group("t.id, t.code, t.name, te.event_id, t.is_active, u.name").
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&teams).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalTeam(&filter)
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

func (repository *TeamRepositoryImpl) getTotalTeam(filter *model.FilterTeam) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterTeam(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	db := repository.DB.Table("teams as t").
		Joins("inner join team_events te on te.team_id = t.id").
		Joins("inner join team_members tm on tm.team_id = t.id").
		Joins("inner join participants p on p.id = t.participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Where(buildWhereQuery, whereVals...)

	if filter.ScheduleID != 0 {
		db.Joins("left join schedule_teams st on st.team_id = t.id")
	}

	if err = db.Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterTeam(filter model.FilterTeam) (where []string, whereVal []interface{}) {
	where = append(where, "t.deleted_at IS NULL")

	if filter.Search != "" {
		filter.Search = strings.ToLower(filter.Search)
		if filter.IsForParticipant {
			where = append(where, "(LOWER(t.name) LIKE @q)")
		} else {
			where = append(where, "(LOWER(t.name) LIKE @q OR LOWER(u.name) LIKE @q)")
		}
		whereVal = append(whereVal, sql.Named("q", "%"+filter.Search+"%"))
	}

	if filter.EventID != 0 {
		where = append(where, "te.event_id = @event")
		whereVal = append(whereVal, sql.Named("event", filter.EventID))
	}

	if filter.Status != "" {
		isActive := false
		if filter.Status == "active" {
			isActive = true
		}
		where = append(where, "t.is_active = @is_active")
		whereVal = append(whereVal, sql.Named("is_active", isActive))
	}

	if filter.ScheduleID != 0 {
		where = append(where, "st.schedule_id = @schedule")
		whereVal = append(whereVal, sql.Named("schedule", filter.ScheduleID))
	}

	return
}

func (repository *TeamRepositoryImpl) FindDetail(id uint) (team model.TeamDetail, err error) {
	// TODO: project id, project link
	if err = repository.DB.Table("teams as t").
		Select(`t.*, te.event_id, count(tm.id) as num_of_member, u.name as participant_name`).
		Joins("inner join team_events te on te.team_id = t.id").
		Joins("inner join team_members tm on tm.team_id = t.id").
		Joins("inner join participants p on p.id = t.participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Group("t.id, t.code, t.name, te.event_id, t.is_active, u.name").
		Where("t.id=?", id).
		First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *TeamRepositoryImpl) FindDetail2(id, eventID, participantID uint, includeMembers bool) (team model.TeamDetail2, err error) {
	qSelect := `t.id, t.code, t.name, t.description, t.avatar, t.participant_id, 
		count(tm.id) as num_of_member, proj.id as project_id`
	if eventID != 0 && participantID != 0 {
		qSelect = fmt.Sprintf(`%s,
			IF((select count(tr.id)
			from team_requests tr
			where tr.team_id = t.id
				and tr.participant_id = %d
				and tr.event_id = %d
				and tr.status='sent') = 0, false, true) as is_requested`, qSelect, participantID, eventID)
	}

	if err = repository.DB.Table("teams as t").
		Select(qSelect).
		Joins("inner join team_members tm on tm.team_id = t.id").
		Joins("left join projects proj on proj.team_id = t.id AND proj.event_id=?").
		Group("t.id, t.code, t.name, t.description, t.avatar, proj.id").
		Where("t.id=?", eventID, id).
		First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}

	if includeMembers {
		if err = repository.DB.Table("team_members as tm").
			Select("tm.id, tm.participant_id, u.name as participant_name, u.avatar as participant_avatar").
			Joins("inner join participants p on p.id = tm.participant_id").
			Joins("inner join users u on u.id = p.user_id").
			Where("tm.team_id=?", id).
			Find(&team.Members).Error; err != nil {
			return
		}
	}

	return
}

func (repository *TeamRepositoryImpl) FindByIDAndEventID(id, eventID uint) (team model.TeamDetail2, err error) {
	if err = repository.DB.Table("teams as t").
		Select(`t.id, t.code, t.name, t.description, t.avatar, t.participant_id, u.name as participant_name,
			u.email as participant_email, count(tm.id) as num_of_member`).
		Joins("inner join team_events te on te.team_id = t.id").
		Joins("inner join team_members tm on tm.team_id = t.id").
		Joins("inner join participants p on p.id = t.participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Group("t.id, t.code, t.name, t.description, t.avatar").
		Where("t.id=? AND te.event_id=?", id, eventID).
		First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}

	return
}
