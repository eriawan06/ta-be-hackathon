package repository

import (
	"be-sagara-hackathon/src/modules/team/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type TeamInvitationRepository interface {
	Create(invitation model.TeamInvitation) error
	CreateBatch(invitations []model.TeamInvitation) error
	Update(id uint, invitation model.TeamInvitation, member *model.TeamMember) error
	Delete(id uint) error
	FindOne(id uint) (invitation model.TeamInvitation, err error)
	FindByCode(code string) (invitation model.TeamInvitation, err error)
	FindAll(
		filter model.FilterInvitation,
		pg *utils.PaginateQueryOffset,
	) (invitations []model.InvitationLite, totalData, totalPage int64, err error)
	FindDetail(id uint) (invitation model.InvitationDetail, err error)
	FindByEventIDAndTeamIDAndParticipantID(eventID, teamID, participantID uint) (invitations []model.TeamInvitation, err error)
	FindManyByTeamIDAndEventID(teamID, eventID uint) (invitations []model.TeamInvitationList, err error)
	FindDetail2(id uint) (invitation model.TeamInvitationDetail, err error)
}

type TeamInvitationRepositoryImpl struct {
	DB *gorm.DB
}

func NewTeamInvitationRepository(db *gorm.DB) TeamInvitationRepository {
	return &TeamInvitationRepositoryImpl{DB: db}
}

func (repository *TeamInvitationRepositoryImpl) Create(invitation model.TeamInvitation) error {
	if err := repository.DB.Create(&invitation).Error; err != nil {
		return err
	}
	return nil
}

func (repository *TeamInvitationRepositoryImpl) CreateBatch(invitations []model.TeamInvitation) error {
	if err := repository.DB.Create(&invitations).Error; err != nil {
		return err
	}
	return nil
}

func (repository *TeamInvitationRepositoryImpl) Update(id uint, invitation model.TeamInvitation, member *model.TeamMember) error {
	tx := repository.DB.Begin()
	if err := tx.Model(&model.TeamInvitation{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"note":       invitation.Note,
			"status":     invitation.Status,
			"proceed_at": invitation.ProceedAt,
			"updated_at": invitation.UpdatedAt,
			"updated_by": invitation.UpdatedBy,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if member != nil {
		if err := tx.Create(&member).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (repository *TeamInvitationRepositoryImpl) Delete(id uint) error {
	if err := repository.DB.Delete(&model.TeamInvitation{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repository *TeamInvitationRepositoryImpl) FindOne(id uint) (invitation model.TeamInvitation, err error) {
	if err = repository.DB.Where("id=?", id).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *TeamInvitationRepositoryImpl) FindByCode(code string) (invitation model.TeamInvitation, err error) {
	if err = repository.DB.Where("code=?", code).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *TeamInvitationRepositoryImpl) FindAll(
	filter model.FilterInvitation,
	pg *utils.PaginateQueryOffset,
) (invitations []model.InvitationLite, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterInvitation(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("team_invitations as ti").
		Select(`ti.id, ti.code, t.id as team_id, t.code as team_code, t.name, 
			t.avatar, ti.status, count(tm.id) as num_of_member`).
		Joins("inner join teams t on t.id = ti.team_id").
		Joins("inner join team_members tm on tm.team_id = t.id").
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Group("t.id, t.code, t.name, t.avatar, ti.id").
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&invitations).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalInvitation(&filter)
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

func (repository *TeamInvitationRepositoryImpl) getTotalInvitation(filter *model.FilterInvitation) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterInvitation(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("team_invitations as ti").
		Joins("inner join teams t on t.id = ti.team_id").
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}

	return totalData, nil
}

func BuildFilterInvitation(filter model.FilterInvitation) (where []string, whereVal []interface{}) {
	where = append(where, "ti.deleted_at IS NULL")

	if filter.EventID != 0 {
		where = append(where, "ti.event_id = @event")
		whereVal = append(whereVal, sql.Named("event", filter.EventID))
	}

	if filter.ParticipantID != 0 {
		where = append(where, "ti.to_participant_id = @participant")
		whereVal = append(whereVal, sql.Named("participant", filter.ParticipantID))
	}

	if filter.Status != "" {
		where = append(where, "ti.status = @status")
		whereVal = append(whereVal, sql.Named("status", filter.Status))
	}

	return
}

func (repository *TeamInvitationRepositoryImpl) FindDetail(id uint) (invitation model.InvitationDetail, err error) {
	if err = repository.DB.Table("team_invitations as ti").
		Select(`ti.id, ti.code, ti.note, ti.to_participant_id, t.id as team_id, t.code as team_code, t.name, 
			t.description, t.avatar, ti.status, count(tm.id) as num_of_member`).
		Joins("inner join teams t on t.id = ti.team_id").
		Joins("inner join team_members tm on tm.team_id = t.id").
		Group("ti.id, ti.code, ti.note, t.id, t.code, t.name, t.description, t.avatar").
		Where("ti.id=?", id).
		First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}

	if err = repository.DB.Table("team_members as tm").
		Select("tm.id, tm.participant_id, u.name as participant_name, u.avatar as participant_avatar").
		Joins("inner join participants p on p.id = tm.participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Joins("left join team_invitations ti on ti.team_id = tm.team_id").
		Where("ti.id=?", id).
		Find(&invitation.TeamMembers).Error; err != nil {
		return
	}

	return
}

func (repository *TeamInvitationRepositoryImpl) FindByEventIDAndTeamIDAndParticipantID(eventID, teamID, participantID uint) (invitations []model.TeamInvitation, err error) {
	if err = repository.DB.Where("event_id=? AND team_id=? AND to_participant_id=?", eventID, teamID, participantID).
		Find(&invitations).Error; err != nil {
		return
	}
	return
}

func (repository *TeamInvitationRepositoryImpl) FindManyByTeamIDAndEventID(teamID, eventID uint) (invitations []model.TeamInvitationList, err error) {
	if err = repository.DB.Table("team_invitations as ti").
		Select(`p.id, p.user_id, u.avatar, u.name, p.speciality_id, spe.name as speciality_name, 
			p.city_id, city.name as city_name, p.gender, u.occupation_id, occ.name as occupation_name, 
			u.institution, p.school, p.bio, ti.id as invitation_id, ti.code as invitation_code, ti.status`).
		Joins("inner join teams t on t.id = ti.team_id").
		Joins("inner join participants p on p.id = ti.to_participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Joins("inner join occupations occ on occ.id = u.occupation_id").
		Joins("inner join specialities spe on spe.id = p.speciality_id").
		Joins("left join reg_cities city on city.id = p.city_id").
		Where("ti.team_id=? AND ti.event_id=?", teamID, eventID).
		Order(`case ti.status
			when 'sent' then 1
			when 'rejected' then 2
			else 3 end`).
		Find(&invitations).Error; err != nil {
		return
	}

	for k, v := range invitations {
		if err = repository.DB.Table("participant_skills ps").
			Select("ps.skill_id, sk.name as skill_name").
			Joins("inner join skills sk on sk.id = ps.skill_id").
			Where("participant_id=?", v.ID).
			Find(&invitations[k].Skills).Error; err != nil {
			return
		}
	}
	return
}

func (repository *TeamInvitationRepositoryImpl) FindDetail2(id uint) (invitation model.TeamInvitationDetail, err error) {
	if err = repository.DB.Table("team_invitations as ti").
		Select(`p.id, p.user_id, u.avatar, u.name, p.speciality_id, spe.name as speciality_name, 
			p.city_id, city.name as city_name, p.gender, u.occupation_id, occ.name as occupation_name, 
			u.institution, p.school, p.bio, ti.id as invitation_id, ti.code as invitation_code, ti.status, 
			ti.created_at, ti.updated_at, ti.proceed_at, ti.note, ti.team_id`).
		Joins("inner join teams t on t.id = ti.team_id").
		Joins("inner join participants p on p.id = ti.to_participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Joins("inner join occupations occ on occ.id = u.occupation_id").
		Joins("inner join specialities spe on spe.id = p.speciality_id").
		Joins("left join reg_cities city on city.id = p.city_id").
		Where("ti.id=?", id).
		First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}

	if err = repository.DB.Table("participant_skills ps").
		Select("ps.skill_id, sk.name as skill_name").
		Joins("inner join skills sk on sk.id = ps.skill_id").
		Where("participant_id=?", invitation.ID).
		Find(&invitation.Skills).Error; err != nil {
		return
	}
	return
}
