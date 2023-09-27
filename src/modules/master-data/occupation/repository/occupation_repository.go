package repository

import (
	"be-sagara-hackathon/src/modules/master-data/occupation/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type OccupationRepository interface {
	Save(req model.Occupation) (err error)
	Update(req model.Occupation, id uint) (err error)
	Find(
		filter model.FilterOccupation,
		pg *utils.PaginateQueryOffset,
	) (occupations []model.Occupation, totalData, totalPage int64, err error)
	FindOne(id uint) (occupation model.Occupation, err error)
}

type OccupationRepositoryImpl struct {
	DB *gorm.DB
}

func NewOccupationRepository(db *gorm.DB) OccupationRepository {
	return &OccupationRepositoryImpl{DB: db}
}

func (repository *OccupationRepositoryImpl) Save(req model.Occupation) (err error) {
	if err = repository.DB.Create(&req).Error; err != nil {
		return
	}
	return
}

func (repository *OccupationRepositoryImpl) Update(req model.Occupation, id uint) (err error) {
	if err = repository.DB.Select("*").Where("id = ?", id).Updates(&req).Error; err != nil {
		return
	}
	return
}

func (repository *OccupationRepositoryImpl) Find(
	filter model.FilterOccupation,
	pg *utils.PaginateQueryOffset,
) (occupations []model.Occupation, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilter(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&occupations).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalOccupation(&filter)
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

func (repository *OccupationRepositoryImpl) getTotalOccupation(filter *model.FilterOccupation) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilter(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.Occupation{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilter(filter model.FilterOccupation) (where []string, whereVal []interface{}) {
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

func (repository *OccupationRepositoryImpl) FindOne(id uint) (occupation model.Occupation, err error) {
	if err = repository.DB.Where("id=?", id).First(&occupation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}
