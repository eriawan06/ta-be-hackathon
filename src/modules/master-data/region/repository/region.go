package repository

import (
	"be-sagara-hackathon/src/modules/master-data/region/model"
	"be-sagara-hackathon/src/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type RegionRepository interface {
	FindProvince(
		ctx context.Context,
		filter model.FilterProvince,
		pg *utils.PaginateQueryOffset,
	) (provinces []model.RegProvince, totalData, totalPage int64, err error)
	FindCity(
		ctx context.Context,
		filter model.FilterCity,
		pg *utils.PaginateQueryOffset,
	) (cities []model.RegCity, totalData, totalPage int64, err error)
	FindDistrict(
		ctx context.Context,
		filter model.FilterDistrict,
		pg *utils.PaginateQueryOffset,
	) (districts []model.RegDistrict, totalData, totalPage int64, err error)
	FindVillage(
		ctx context.Context,
		filter model.FilterVillage,
		pg *utils.PaginateQueryOffset,
	) (villages []model.RegVillage, totalData, totalPage int64, err error)
}

type RegionRepositoryImpl struct {
	db *gorm.DB
}

func NewRegionRepository(database *gorm.DB) RegionRepository {
	return &RegionRepositoryImpl{db: database}
}

func (repository *RegionRepositoryImpl) FindProvince(
	ctx context.Context,
	filter model.FilterProvince,
	pg *utils.PaginateQueryOffset,
) (provinces []model.RegProvince, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterProvince(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.db.WithContext(ctx).
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&provinces).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalProvince(&filter)
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

func (repository *RegionRepositoryImpl) getTotalProvince(filter *model.FilterProvince) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterProvince(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.db.Model(&model.RegProvince{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterProvince(filter model.FilterProvince) (where []string, whereVal []interface{}) {
	if filter.ID != "" {
		where = append(where, "id = ?")
		whereVal = append(whereVal, filter.ID)
	}

	if filter.Name != "" {
		where = append(where, "LOWER(name) LIKE ?")
		whereVal = append(whereVal, "%"+strings.ToLower(filter.Name)+"%")
	}

	return
}

func (repository *RegionRepositoryImpl) FindCity(
	ctx context.Context,
	filter model.FilterCity,
	pg *utils.PaginateQueryOffset,
) (cities []model.RegCity, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterCity(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.db.WithContext(ctx).
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&cities).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalCity(&filter)
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

func (repository *RegionRepositoryImpl) getTotalCity(filter *model.FilterCity) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterCity(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.db.Model(&model.RegCity{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterCity(filter model.FilterCity) (where []string, whereVal []interface{}) {

	if filter.ID != "" {
		where = append(where, "id = ?")
		whereVal = append(whereVal, filter.ID)
	}

	if filter.ProvinceID != "" {
		where = append(where, "province_id = ?")
		whereVal = append(whereVal, filter.ProvinceID)
	}

	if filter.Name != "" {
		where = append(where, "LOWER(name) LIKE ?")
		whereVal = append(whereVal, "%"+strings.ToLower(filter.Name)+"%")
	}

	return
}

func (repository *RegionRepositoryImpl) FindDistrict(
	ctx context.Context,
	filter model.FilterDistrict,
	pg *utils.PaginateQueryOffset,
) (districts []model.RegDistrict, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterDistrict(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.db.WithContext(ctx).
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&districts).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalDistrict(&filter)
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

func (repository *RegionRepositoryImpl) getTotalDistrict(filter *model.FilterDistrict) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterDistrict(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.db.Model(&model.RegDistrict{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterDistrict(filter model.FilterDistrict) (where []string, whereVal []interface{}) {
	if filter.ID != "" {
		where = append(where, "id = ?")
		whereVal = append(whereVal, filter.ID)
	}

	if filter.CityID != "" {
		where = append(where, "city_id = ?")
		whereVal = append(whereVal, filter.CityID)
	}

	if filter.Name != "" {
		where = append(where, "LOWER(name) LIKE ?")
		whereVal = append(whereVal, "%"+strings.ToLower(filter.Name)+"%")
	}

	return
}

func (repository *RegionRepositoryImpl) FindVillage(
	ctx context.Context,
	filter model.FilterVillage,
	pg *utils.PaginateQueryOffset,
) (villages []model.RegVillage, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterVillage(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.db.WithContext(ctx).
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&villages).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalVillage(&filter)
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

func (repository *RegionRepositoryImpl) getTotalVillage(filter *model.FilterVillage) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterVillage(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.db.Model(&model.RegVillage{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterVillage(filter model.FilterVillage) (where []string, whereVal []interface{}) {
	if filter.ID != "" {
		where = append(where, "id = ?")
		whereVal = append(whereVal, filter.ID)
	}

	if filter.DistrictID != "" {
		where = append(where, "district_id = ?")
		whereVal = append(whereVal, filter.DistrictID)
	}

	if filter.Name != "" {
		where = append(where, "LOWER(name) LIKE ?")
		whereVal = append(whereVal, "%"+strings.ToLower(filter.Name)+"%")
	}

	return
}
