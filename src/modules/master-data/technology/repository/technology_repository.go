package repository

import (
	"be-sagara-hackathon/src/modules/master-data/technology/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type TechnologyRepository interface {
	Save(req model.Technology) (err error)
	Update(req model.Technology, id uint) (err error)
	Find(
		filter model.FilterTechnology,
		pg *utils.PaginateQueryOffset,
	) (technologies []model.Technology, totalData, totalPage int64, err error)
	FindOne(id uint) (technology model.Technology, err error)
}

type TechnologyRepositoryImpl struct {
	DB *gorm.DB
}

func NewTechnologyRepository(db *gorm.DB) TechnologyRepository {
	return &TechnologyRepositoryImpl{DB: db}
}

func (repository *TechnologyRepositoryImpl) Save(req model.Technology) (err error) {
	if err = repository.DB.Create(&req).Error; err != nil {
		return
	}
	return
}

func (repository *TechnologyRepositoryImpl) Update(req model.Technology, id uint) (err error) {
	if err = repository.DB.Select("*").Where("id = ?", id).Updates(&req).Error; err != nil {
		return
	}
	return
}

func (repository *TechnologyRepositoryImpl) Find(
	filter model.FilterTechnology,
	pg *utils.PaginateQueryOffset,
) (technologies []model.Technology, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilter(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&technologies).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalTechnology(&filter)
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

func (repository *TechnologyRepositoryImpl) getTotalTechnology(filter *model.FilterTechnology) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilter(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.Technology{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilter(filter model.FilterTechnology) (where []string, whereVal []interface{}) {
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

func (repository *TechnologyRepositoryImpl) FindOne(id uint) (technology model.Technology, err error) {
	if err = repository.DB.Where("id=?", id).First(&technology).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}
