package service

import (
	"be-sagara-hackathon/src/modules/team/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	e "be-sagara-hackathon/src/utils/errors"
	"context"
)

type TeamMemberService interface {
	Delete(ctx context.Context, id uint) error
}

type TeamMemberServiceImpl struct {
	Repository repository.TeamMemberRepository
	//TeamRepository repository.TeamRepository
}

func NewTeamMemberService(
	repository repository.TeamMemberRepository,
	// teamRepository repository.TeamRepository,
) TeamMemberService {
	return &TeamMemberServiceImpl{
		Repository: repository,
		//TeamRepository: teamRepository,
	}
}

func (service *TeamMemberServiceImpl) Delete(ctx context.Context, id uint) error {
	authenticatedUser := ctx.Value("user").(um.User)
	teamMember, err := service.Repository.FindByID(id)
	if err != nil {
		return err
	}

	if teamMember.Team.ParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	//can not remove team's admin/creator
	if teamMember.Team.ParticipantID == teamMember.ParticipantID {
		return e.ErrCannotRemoveTeamAdmin
	}

	if err = service.Repository.Delete(id); err != nil {
		return err
	}

	return nil
}
