package controller

import (
	"be-sagara-hackathon/src/modules/master-data/region/model"
	"be-sagara-hackathon/src/modules/master-data/region/service"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegionController interface {
	GetListProvince(ctx *gin.Context)
	GetListCity(ctx *gin.Context)
	GetListDistrict(ctx *gin.Context)
	GetListVillage(ctx *gin.Context)
}

type RegionControllerImpl struct {
	Service service.RegionService
}

func NewRegionController(service service.RegionService) RegionController {
	return &RegionControllerImpl{Service: service}
}

func (controller *RegionControllerImpl) GetListProvince(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	var filter model.FilterProvince
	filter.ID = ctx.Query("id")
	filter.Name = ctx.Query("name")

	data, err := controller.Service.GetListProvince(ctx, filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Province Success", data)
}

func (controller *RegionControllerImpl) GetListCity(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	var filter model.FilterCity
	filter.ID = ctx.Query("id")
	filter.ProvinceID = ctx.Query("province")
	filter.Name = ctx.Query("name")

	data, err := controller.Service.GetListCity(ctx, filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List City Success", data)
}

func (controller *RegionControllerImpl) GetListDistrict(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	var filter model.FilterDistrict
	filter.ID = ctx.Query("id")
	filter.CityID = ctx.Query("city")
	filter.Name = ctx.Query("name")

	data, err := controller.Service.GetListDistrict(ctx, filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List District Success", data)
}

func (controller *RegionControllerImpl) GetListVillage(ctx *gin.Context) {
	pg, err := utils.GetPaginateQueryOffset(ctx.Request)
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Not Found Error", []string{err.Error()})
		return
	}

	var filter model.FilterVillage
	filter.ID = ctx.Query("id")
	filter.DistrictID = ctx.Query("district")
	filter.Name = ctx.Query("name")

	data, err := controller.Service.GetListVillage(ctx, filter, pg)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get List Village Success", data)
}
