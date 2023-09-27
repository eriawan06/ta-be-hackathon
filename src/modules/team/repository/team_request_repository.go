package repository

import (
	"be-sagara-hackathon/src/modules/team/model"
	e "be-sagara-hackathon/src/utils/errors"
	"gorm.io/gorm"
)

type TeamRequestRepository interface {
	Create(request model.TeamRequest) error
	Update(id uint, request model.TeamRequest, member *model.TeamMember) error
	Delete(id uint) error
	FindOne(id uint) (teamReq model.TeamRequest, err error)
	FindByCode(code string) (teamReq model.TeamRequest, err error)
	FindByEventIDAndTeamIDAndParticipantID(eventID, teamID, participantID uint) (teamReqs []model.TeamRequest, err error)
	FindManyByTeamIDAndEventID(teamID, eventID uint) (teamReqs []model.TeamRequestList, err error)
	FindDetail(id uint) (teamReq model.TeamRequestDetail, err error)
}

type TeamRequestRepositoryImpl struct {
	DB *gorm.DB
}

func NewTeamRequestRepository(db *gorm.DB) TeamRequestRepository {
	return &TeamRequestRepositoryImpl{DB: db}
}

func (repository *TeamRequestRepositoryImpl) Create(request model.TeamRequest) error {
	if err := repository.DB.Create(&request).Error; err != nil {
		return err
	}
	return nil
}

func (repository *TeamRequestRepositoryImpl) Update(id uint, request model.TeamRequest, member *model.TeamMember) error {
	tx := repository.DB.Begin()
	if err := tx.Model(&model.TeamRequest{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"note":       request.Note,
			"status":     request.Status,
			"proceed_at": request.ProceedAt,
			"proceed_by": request.ProceedBy,
			"updated_at": request.UpdatedAt,
			"updated_by": request.UpdatedBy,
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

func (repository *TeamRequestRepositoryImpl) Delete(requestId uint) error {
	if err := repository.DB.Delete(&model.TeamRequest{}, requestId).Error; err != nil {
		return err
	}
	return nil
}

func (repository *TeamRequestRepositoryImpl) FindOne(id uint) (teamReq model.TeamRequest, err error) {
	if err = repository.DB.Where("id=?", id).
		Preload("Team").
		First(&teamReq).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *TeamRequestRepositoryImpl) FindByCode(code string) (teamReq model.TeamRequest, err error) {
	if err = repository.DB.Where("code=?", code).First(&teamReq).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *TeamRequestRepositoryImpl) FindByEventIDAndTeamIDAndParticipantID(eventID, teamID, participantID uint) (teamReqs []model.TeamRequest, err error) {
	if err = repository.DB.Where("event_id=? AND team_id=? AND participant_id=?", eventID, teamID, participantID).
		Find(&teamReqs).Error; err != nil {
		return
	}
	return
}

func (repository *TeamRequestRepositoryImpl) FindManyByTeamIDAndEventID(teamID, eventID uint) (teamReqs []model.TeamRequestList, err error) {
	if err = repository.DB.Table("team_requests as tr").
		Select(`p.id, p.user_id, u.avatar, u.name, p.speciality_id, spe.name as speciality_name, 
			p.city_id, city.name as city_name, p.gender, u.occupation_id, occ.name as occupation_name, 
			u.institution, p.school, p.bio, tr.id as request_id, tr.code as request_code, tr.status`).
		Joins("inner join teams t on t.id = tr.team_id").
		Joins("inner join participants p on p.id = tr.participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Joins("inner join occupations occ on occ.id = u.occupation_id").
		Joins("inner join specialities spe on spe.id = p.speciality_id").
		Joins("left join reg_cities city on city.id = p.city_id").
		Where("tr.team_id=? AND tr.event_id=?", teamID, eventID).
		Order(`case tr.status
			when 'sent' then 1
			when 'rejected' then 2
			else 3 end`).
		Find(&teamReqs).Error; err != nil {
		return
	}

	for k, v := range teamReqs {
		if err = repository.DB.Table("participant_skills ps").
			Select("ps.skill_id, sk.name as skill_name").
			Joins("inner join skills sk on sk.id = ps.skill_id").
			Where("participant_id=?", v.ID).
			Find(&teamReqs[k].Skills).Error; err != nil {
			return
		}
	}
	return
}

func (repository *TeamRequestRepositoryImpl) FindDetail(id uint) (teamReq model.TeamRequestDetail, err error) {
	if err = repository.DB.Table("team_requests as tr").
		Select(`p.id, p.user_id, u.avatar, u.name, p.speciality_id, spe.name as speciality_name, 
			p.city_id, city.name as city_name, p.gender, u.occupation_id, occ.name as occupation_name, 
			u.institution, p.school, p.bio, tr.id as request_id, tr.code as request_code, tr.status, tr.created_at,
			tr.updated_at, tr.proceed_at, tr.proceed_by, tr.note, tr.team_id`).
		Joins("inner join teams t on t.id = tr.team_id").
		Joins("inner join participants p on p.id = tr.participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Joins("inner join occupations occ on occ.id = u.occupation_id").
		Joins("inner join specialities spe on spe.id = p.speciality_id").
		Joins("left join reg_cities city on city.id = p.city_id").
		Where("tr.id=?", id).
		First(&teamReq).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}

	if err = repository.DB.Table("participant_skills ps").
		Select("ps.skill_id, sk.name as skill_name").
		Joins("inner join skills sk on sk.id = ps.skill_id").
		Where("participant_id=?", teamReq.ID).
		Find(&teamReq.Skills).Error; err != nil {
		return
	}
	return
}
