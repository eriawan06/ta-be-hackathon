package repository

import (
	"be-sagara-hackathon/src/modules/master-data/speciality/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type SpecialityRepository interface {
	Save(req model.Speciality) (err error)
	Update(req model.Speciality, id uint) (err error)
	Find(
		filter model.FilterSpeciality,
		pg *utils.PaginateQueryOffset,
	) (specialities []model.Speciality, totalData, totalPage int64, err error)
	FindOne(id uint) (speciality model.Speciality, err error)
}

type SpecialityRepositoryImpl struct {
	DB *gorm.DB
}

func NewSpecialityRepository(db *gorm.DB) SpecialityRepository {
	return &SpecialityRepositoryImpl{DB: db}
}

func (repository *SpecialityRepositoryImpl) Save(req model.Speciality) (err error) {
	if err = repository.DB.Create(&req).Error; err != nil {
		return
	}
	return
}

func (repository *SpecialityRepositoryImpl) Update(req model.Speciality, id uint) (err error) {
	if err = repository.DB.Select("*").Where("id = ?", id).Updates(&req).Error; err != nil {
		return
	}
	return
}

func (repository *SpecialityRepositoryImpl) Find(
	filter model.FilterSpeciality,
	pg *utils.PaginateQueryOffset,
) (specialities []model.Speciality, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilter(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&specialities).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalSpeciality(&filter)
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

func (repository *SpecialityRepositoryImpl) getTotalSpeciality(filter *model.FilterSpeciality) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilter(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.Speciality{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilter(filter model.FilterSpeciality) (where []string, whereVal []interface{}) {
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

func (repository *SpecialityRepositoryImpl) FindOne(id uint) (speciality model.Speciality, err error) {
	if err = repository.DB.Where("id=?", id).First(&speciality).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}
