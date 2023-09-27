package router

import (
	"be-sagara-hackathon/src/modules/master-data/region"
	"github.com/gin-gonic/gin"
)

func RegionRouter(group *gin.RouterGroup) {
	group.GET("/provinces", region.GetController().GetListProvince)
	group.GET("/cities", region.GetController().GetListCity)
	group.GET("/districts", region.GetController().GetListDistrict)
	group.GET("/villages", region.GetController().GetListVillage)
}
