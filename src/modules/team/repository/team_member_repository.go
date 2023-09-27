package repository

import (
	"be-sagara-hackathon/src/modules/team/model"
	e "be-sagara-hackathon/src/utils/errors"
	"errors"
	"gorm.io/gorm"
)

type TeamMemberRepository interface {
	Save(member model.TeamMember) error
	SaveBatch(members []model.TeamMember) error
	Delete(id uint) error
	FindByID(id uint) (member model.TeamMember, err error)
	FindByParticipantID(participantID uint) (member model.TeamMember, err error)
	FindByParticipantIDAndTeamID(participantID, teamID uint) (member model.TeamMember, err error)
	FindManyByTeamID(teamID uint) (members []model.TeamMemberList, err error)
	//FindByTeamIDAndParticipantID(teamID, participantID uint) (model.TeamMemberDetail, error)
	//FindByTeamCodeAndParticipantID(teamCode string, participantID uint) (model.TeamMemberDetail, error)
}

type TeamMemberRepositoryImpl struct {
	DB *gorm.DB
}

func NewTeamMemberRepository(db *gorm.DB) TeamMemberRepository {
	return &TeamMemberRepositoryImpl{DB: db}
}

func (repository *TeamMemberRepositoryImpl) Save(member model.TeamMember) error {
	if err := repository.DB.Create(&member).Error; err != nil {
		return err
	}
	return nil
}

func (repository *TeamMemberRepositoryImpl) SaveBatch(members []model.TeamMember) error {
	if err := repository.DB.Create(&members).Error; err != nil {
		return err
	}
	return nil
}

func (repository *TeamMemberRepositoryImpl) Delete(id uint) error {
	if err := repository.DB.Delete(&model.TeamMember{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repository *TeamMemberRepositoryImpl) FindByID(id uint) (member model.TeamMember, err error) {
	if err = repository.DB.Where("id=?", id).
		Preload("Team").
		First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *TeamMemberRepositoryImpl) FindByParticipantID(participantID uint) (member model.TeamMember, err error) {
	if err = repository.DB.Where("participant_id=?", participantID).
		Preload("Team").
		First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *TeamMemberRepositoryImpl) FindByParticipantIDAndTeamID(participantID, teamID uint) (member model.TeamMember, err error) {
	if err = repository.DB.Where("participant_id=? AND team_id=?", participantID, teamID).
		First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *TeamMemberRepositoryImpl) FindManyByTeamID(teamID uint) (members []model.TeamMemberList, err error) {
	if err = repository.DB.Table("team_members as tm").
		Select(`p.id, p.user_id, u.avatar, u.name, p.speciality_id, spe.name as speciality_name, 
			p.city_id, city.name as city_name, p.gender, u.occupation_id, occ.name as occupation_name, 
			u.institution, p.school, p.bio, tm.id as team_member_id, tm.joined_at,
			if(t.participant_id = p.id, true, false) as is_admin`).
		Joins("inner join teams t on t.id = tm.team_id").
		Joins("inner join participants p on p.id = tm.participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Joins("inner join occupations occ on occ.id = u.occupation_id").
		Joins("inner join specialities spe on spe.id = p.speciality_id").
		Joins("left join reg_cities city on city.id = p.city_id").
		Where("tm.team_id=?", teamID).
		Find(&members).Error; err != nil {
		return
	}

	for k, v := range members {
		if err = repository.DB.Table("participant_skills ps").
			Select("ps.skill_id, sk.name as skill_name").
			Joins("inner join skills sk on sk.id = ps.skill_id").
			Where("participant_id=?", v.ID).
			Find(&members[k].Skills).Error; err != nil {
			return
		}
	}
	return
}
