package controller

import (
	"be-sagara-hackathon/src/modules/home/service"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HomeController interface {
	GetData(ctx *gin.Context)
}

type HomeControllerImpl struct {
	Service service.HomeService
}

func NewHomeController(service service.HomeService) HomeController {
	return &HomeControllerImpl{Service: service}
}

// GetData Get Home Data godoc
// @Tags Home
// @Summary Get Home Data
// @Description Get Home Data
// @Produce  json
// @Success 200 {object} src.GetEventSuccess
// @Success 400 {object} src.BaseFailure
// @Router /home [get]
func (controller HomeControllerImpl) GetData(ctx *gin.Context) {
	data, err := controller.Service.GetData()

	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendSuccess(ctx, http.StatusOK, "Get Home Data Success", nil)
			return
		}
		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Get Home Data Success", data)
}
