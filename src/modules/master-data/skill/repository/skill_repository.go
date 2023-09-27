package repository

import (
	"be-sagara-hackathon/src/modules/master-data/skill/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type SkillRepository interface {
	Save(req model.Skill) (err error)
	Update(req model.Skill, id uint) (err error)
	Find(
		filter model.FilterSkill,
		pg *utils.PaginateQueryOffset,
	) (skills []model.Skill, totalData, totalPage int64, err error)
	FindOne(id uint) (skill model.Skill, err error)
}

type SkillRepositoryImpl struct {
	DB *gorm.DB
}

func NewSkillRepository(db *gorm.DB) SkillRepository {
	return &SkillRepositoryImpl{DB: db}
}

func (repository *SkillRepositoryImpl) Save(req model.Skill) (err error) {
	if err = repository.DB.Create(&req).Error; err != nil {
		return
	}
	return
}

func (repository *SkillRepositoryImpl) Update(req model.Skill, id uint) (err error) {
	if err = repository.DB.Select("*").Where("id = ?", id).Updates(&req).Error; err != nil {
		return
	}
	return
}

func (repository *SkillRepositoryImpl) Find(
	filter model.FilterSkill,
	pg *utils.PaginateQueryOffset,
) (skills []model.Skill, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilter(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&skills).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalSkill(&filter)
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

func (repository *SkillRepositoryImpl) getTotalSkill(filter *model.FilterSkill) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilter(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.Skill{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilter(filter model.FilterSkill) (where []string, whereVal []interface{}) {
	if filter.Name != "" {
		where = append(where, "LOWER(name) LIKE ?")
		whereVal = append(whereVal, "%"+filter.Name+"%")
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

func (repository *SkillRepositoryImpl) FindOne(id uint) (skill model.Skill, err error) {
	if err = repository.DB.Where("id=?", id).First(&skill).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}
