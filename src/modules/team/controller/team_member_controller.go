package controller

import (
	"be-sagara-hackathon/src/modules/team/service"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TeamMemberController interface {
	Delete(ctx *gin.Context)
}

type TeamMemberControllerImpl struct {
	Service service.TeamMemberService
}

func NewTeamMemberController(service service.TeamMemberService) TeamMemberController {
	return &TeamMemberControllerImpl{Service: service}
}

func (controller *TeamMemberControllerImpl) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid Id", []string{err.Error()})
		return
	}

	err = controller.Service.Delete(ctx, uint(id))
	if err != nil {
		if err == e.ErrDataNotFound {
			common.SendError(ctx, http.StatusNotFound, "Not Found", []string{err.Error()})
			return
		}

		if err == e.ErrForbidden {
			common.SendError(ctx, http.StatusForbidden, "Forbidden", []string{err.Error()})
			return
		}

		if err == e.ErrCannotRemoveTeamAdmin {
			common.SendError(ctx, http.StatusBadRequest, "Bad Request", []string{err.Error()})
			return
		}

		common.SendError(ctx, http.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return
	}

	common.SendSuccess(ctx, http.StatusOK, "Delete Team Member Success", nil)
}
