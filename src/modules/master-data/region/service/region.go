package service

import (
	"be-sagara-hackathon/src/modules/master-data/region/model"
	"be-sagara-hackathon/src/modules/master-data/region/repository"
	"be-sagara-hackathon/src/utils"
	"context"
)

type RegionService interface {
	GetListProvince(
		ctx context.Context,
		filter model.FilterProvince,
		pg *utils.PaginateQueryOffset,
	) (response model.ListProvinceResponse, err error)
	GetListCity(
		ctx context.Context,
		filter model.FilterCity,
		pg *utils.PaginateQueryOffset,
	) (response model.ListCityResponse, err error)
	GetListDistrict(
		ctx context.Context,
		filter model.FilterDistrict,
		pg *utils.PaginateQueryOffset,
	) (response model.ListDistrictResponse, err error)
	GetListVillage(
		ctx context.Context,
		filter model.FilterVillage,
		pg *utils.PaginateQueryOffset,
	) (response model.ListVillageResponse, err error)
}

type RegionServiceImpl struct {
	Repository repository.RegionRepository
}

func NewRegionService(repository repository.RegionRepository) RegionService {
	return &RegionServiceImpl{Repository: repository}
}

func (service *RegionServiceImpl) GetListProvince(
	ctx context.Context,
	filter model.FilterProvince,
	pg *utils.PaginateQueryOffset,
) (response model.ListProvinceResponse, err error) {
	provinces, totalData, totalPage, err := service.Repository.FindProvince(ctx, filter, pg)
	if err != nil {
		return
	}
	response.Provinces = provinces
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *RegionServiceImpl) GetListCity(
	ctx context.Context,
	filter model.FilterCity,
	pg *utils.PaginateQueryOffset,
) (response model.ListCityResponse, err error) {
	cities, totalData, totalPage, err := service.Repository.FindCity(ctx, filter, pg)
	if err != nil {
		return
	}
	response.Cities = cities
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *RegionServiceImpl) GetListDistrict(
	ctx context.Context,
	filter model.FilterDistrict,
	pg *utils.PaginateQueryOffset,
) (response model.ListDistrictResponse, err error) {
	districts, totalData, totalPage, err := service.Repository.FindDistrict(ctx, filter, pg)
	if err != nil {
		return
	}
	response.Districts = districts
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *RegionServiceImpl) GetListVillage(
	ctx context.Context,
	filter model.FilterVillage,
	pg *utils.PaginateQueryOffset,
) (response model.ListVillageResponse, err error) {
	villages, totalData, totalPage, err := service.Repository.FindVillage(ctx, filter, pg)
	if err != nil {
		return
	}
	response.Villages = villages
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}
