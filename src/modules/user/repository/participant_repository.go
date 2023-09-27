package repository

import (
	evm "be-sagara-hackathon/src/modules/event/model"
	pym "be-sagara-hackathon/src/modules/payment/model"
	"be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type ParticipantRepository interface {
	Save(participant model.Participant) (err error)
	Update(req model.UpdateParticipant) (err error)
	UpdateRegistrationStatus(
		id uint,
		isRegistered bool,
		eventParticipant *evm.EventParticipant,
		invoice *pym.Invoice,
	) (err error)
	FindByID(id uint) (participant model.Participant, err error)
	FindByIDs(ids []uint) (participants []model.Participant, err error)
	FindByEmail(email string) (participant model.Participant, err error)
	FindDetail(id uint) (participant model.Participant, err error)
	Find(
		filter model.FilterParticipantSearch,
		pg *utils.PaginateQueryOffset,
	) (participants []model.ParticipantSearch, totalData, totalPage int64, err error)
}

type ParticipantRepositoryImpl struct {
	DB *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) ParticipantRepository {
	return &ParticipantRepositoryImpl{DB: db}
}

func (repository *ParticipantRepositoryImpl) Save(participant model.Participant) (err error) {
	participant.BaseEntity.CreatedBy = "self"
	participant.BaseEntity.UpdatedBy = "self"
	result := repository.DB.Create(&participant)

	return result.Error
}

func (repository *ParticipantRepositoryImpl) Update(req model.UpdateParticipant) (err error) {
	tx := repository.DB.Begin()
	if err = tx.Model(&model.Participant{}).Where("id = ?", req.ID).Updates(&req.Participant).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Model(&model.User{}).Where("id = ?", req.Participant.UserID).Updates(&req.Participant.User).Error; err != nil {
		tx.Rollback()
		return
	}

	if len(req.Skills) > 0 {
		if err = tx.Create(&req.Skills).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	if len(req.RemovedSkills) > 0 {
		err = tx.Delete(&model.ParticipantSkill{}, "participant_id=? AND skill_id in (?)", req.ID, req.RemovedSkills).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	return
}

func (repository *ParticipantRepositoryImpl) UpdateRegistrationStatus(
	id uint,
	isRegistered bool,
	eventParticipant *evm.EventParticipant,
	invoice *pym.Invoice,
) (err error) {
	tx := repository.DB.Begin()

	updateModel := map[string]interface{}{"is_registered": isRegistered}
	if invoice != nil {
		updateModel["payment_status"] = constants.InvoiceUnpaid
	}

	if err = tx.Model(&model.Participant{}).
		Where("id = ?", id).
		Updates(updateModel).Error; err != nil {
		tx.Rollback()
		return
	}

	if eventParticipant != nil {
		if err = tx.Create(&eventParticipant).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	if invoice != nil {
		if err = tx.Create(&invoice).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	return
}

func (repository *ParticipantRepositoryImpl) FindByID(id uint) (participant model.Participant, err error) {
	if err = repository.DB.Where("id=?", id).First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *ParticipantRepositoryImpl) FindByIDs(ids []uint) (participants []model.Participant, err error) {
	if err = repository.DB.Where("id in (?)", ids).First(&participants).Error; err != nil {
		return
	}
	return
}

func (repository *ParticipantRepositoryImpl) FindByEmail(email string) (participant model.Participant, err error) {
	if err = repository.DB.Joins("inner join users u on u.id = participants.user_id").
		Preload("User").Preload("User.Occupation").
		Preload("Province").Preload("City").Preload("District").Preload("Village").
		Preload("Speciality").
		Preload("Skills").Preload("Skills.Skill").
		Where("u.email=?", email).
		First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *ParticipantRepositoryImpl) FindDetail(id uint) (participant model.Participant, err error) {
	if err = repository.DB.Preload("User").Preload("User.Occupation").
		Preload("Province").Preload("City").Preload("District").Preload("Village").
		Preload("Speciality").
		Preload("Skills").Preload("Skills.Skill").
		Where("id=?", id).
		First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *ParticipantRepositoryImpl) Find(
	filter model.FilterParticipantSearch,
	pg *utils.PaginateQueryOffset,
) (participants []model.ParticipantSearch, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterParticipant(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("participants as p").
		Select(`p.id, p.user_id, u.avatar, u.name, p.speciality_id, spe.name as speciality_name, 
			p.city_id, city.name as city_name, p.gender, u.occupation_id, occ.name as occupation_name, 
			u.institution, p.school, p.bio`).
		Joins("inner join users u on u.id = p.user_id").
		Joins("inner join occupations occ on occ.id = u.occupation_id").
		Joins("inner join specialities spe on spe.id = p.speciality_id").
		Joins("left join reg_cities city on city.id = p.city_id").
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&participants).Error; err != nil {
		return
	}

	for k, v := range participants {
		if err = repository.DB.Table("participant_skills ps").
			Select("ps.skill_id, sk.name as skill_name").
			Joins("inner join skills sk on sk.id = ps.skill_id").
			Where("participant_id=?", v.ID).
			Find(&participants[k].Skills).Error; err != nil {
			return
		}
	}

	totalData, err = repository.getTotalParticipant(&filter)
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

func (repository *ParticipantRepositoryImpl) getTotalParticipant(filter *model.FilterParticipantSearch) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterParticipant(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("participants as p").
		Joins("inner join users u on u.id = p.user_id").
		Joins("inner join occupations occ on occ.id = u.occupation_id").
		Joins("inner join specialities spe on spe.id = p.speciality_id").
		Joins("left join reg_cities city on city.id = p.city_id").
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterParticipant(filter model.FilterParticipantSearch) (where []string, whereVal []interface{}) {
	where = append(where, "p.deleted_at is null AND u.deleted_at is null")

	if filter.Search != "" {
		filter.Search = strings.ToLower(filter.Search)
		where = append(where, "LOWER(u.name) LIKE @keyword")
		whereVal = append(whereVal, sql.Named("keyword", "%"+filter.Search+"%"))
	}

	if !filter.InTeam {
		where = append(where, `(select count(*) from team_members tm where tm.participant_id = p.id) = 0`)
	}

	if len(filter.Specialities) > 0 {
		where = append(where, "spe.id in @specialities")
		whereVal = append(whereVal, sql.Named("specialities", filter.Specialities))
	}

	if len(filter.Skills) > 0 {
		where = append(where, `(select count(*)
        from participant_skills ps
        inner join skills sk on sk.id = ps.skill_id
        where ps.participant_id = p.id and ps.skill_id in @skills) > 0`)
		whereVal = append(whereVal, sql.Named("skills", filter.Skills))
	}

	return
}
