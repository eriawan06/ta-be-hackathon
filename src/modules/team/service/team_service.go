package service

import (
	ever "be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/team/model"
	"be-sagara-hackathon/src/modules/team/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	ur "be-sagara-hackathon/src/modules/user/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"context"
	"time"
)

type TeamService interface {
	Create(ctx context.Context, request model.CreateTeamRequest) (team model.Team, err error)
	Update(ctx context.Context, id uint, request model.UpdateTeamRequest) error
	UpdateStatus(ctx context.Context, id uint, request model.UpdateTeamStatusRequest) error
	Delete(ctx context.Context, id uint) error
	GetAll(
		filter model.FilterTeam,
		pg *utils.PaginateQueryOffset,
	) (response model.GetAllTeamResponse, err error)
	GetDetail(id uint) (team model.TeamDetail, err error)
	GetListByEventID(
		ctx context.Context,
		filter model.FilterTeam,
		pg *utils.PaginateQueryOffset,
	) (response model.GetListTeamByEventIDResponse, err error)
	GetDetail2(ctx context.Context, id uint) (team model.TeamDetail2, err error)
	GetMyTeam(ctx context.Context) (team model.TeamDetail2, err error)
	GetMembers(ctx context.Context, id uint) (members []model.TeamMemberList, err error)
	GetInvitations(ctx context.Context, id uint) (invitations []model.TeamInvitationList, err error)
	GetRequests(ctx context.Context, id uint) (requests []model.TeamRequestList, err error)
}

type TeamServiceImpl struct {
	Repository           repository.TeamRepository
	TeamMemberRepo       repository.TeamMemberRepository
	TeamInvitationRepo   repository.TeamInvitationRepository
	TeamRequestRepo      repository.TeamRequestRepository
	ParticipantRepo      ur.ParticipantRepository
	EventRepo            ever.EventRepository
	EventParticipantRepo ever.EventParticipantRepository
}

func NewTeamService(
	teamRepository repository.TeamRepository,
	memberRepository repository.TeamMemberRepository,
	invitationRepository repository.TeamInvitationRepository,
	teamRequestRepo repository.TeamRequestRepository,
	participantRepository ur.ParticipantRepository,
	eventRepository ever.EventRepository,
	eventParticipantRepo ever.EventParticipantRepository,
) TeamService {
	return &TeamServiceImpl{
		Repository:           teamRepository,
		TeamMemberRepo:       memberRepository,
		TeamInvitationRepo:   invitationRepository,
		TeamRequestRepo:      teamRequestRepo,
		ParticipantRepo:      participantRepository,
		EventRepo:            eventRepository,
		EventParticipantRepo: eventParticipantRepo,
	}
}

func (service *TeamServiceImpl) Create(ctx context.Context, request model.CreateTeamRequest) (team model.Team, err error) {
	authenticatedUser := ctx.Value("user").(um.User)

	//get related event
	event, err := service.EventRepo.FindOne(request.EventID)
	if err != nil {
		return
	}

	if event.Status != constants.EventRunning {
		err = e.ErrEventNotRunning
		return
	}

	//get participant (creator)
	participant, err := service.ParticipantRepo.FindByEmail(authenticatedUser.Email)
	if err != nil {
		return
	}

	if !participant.IsRegistered {
		err = e.ErrNotCompleteProfile
		return
	}

	if participant.PaymentStatus != constants.InvoicePaid {
		err = e.ErrPaymentNotPaid
		return
	}

	//check participant has team?
	member, err := service.TeamMemberRepo.FindByParticipantID(participant.ID)
	if err != nil && err != e.ErrDataNotFound {
		return
	}
	if member.ID != 0 && member.Team.DeletedAt == nil && member.Team.IsActive {
		err = e.ErrHasTeam
		return
	}

	// create team
	teamCode := utils.GenerateUuid()
	team = model.Team{
		BaseEntity:    builder.BuildBaseEntity(ctx, true, nil),
		ParticipantID: participant.ID,
		Code:          teamCode,
		Name:          request.Name,
		Description:   request.Description,
		Avatar:        request.Avatar,
	}
	teamMember := model.TeamMember{
		BaseEntity:    builder.BuildBaseEntity(ctx, true, nil),
		ParticipantID: participant.ID,
		JoinedAt:      time.Now(),
	}
	if team, err = service.Repository.Save(request.EventID, team, teamMember); err != nil {
		return
	}

	return
}

func (service *TeamServiceImpl) Update(ctx context.Context, id uint, request model.UpdateTeamRequest) error {
	authenticatedUser := ctx.Value("user").(um.User)
	team, err := service.Repository.FindOne(id)
	if err != nil {
		return err
	}

	if team.ParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	team.Name = request.Name
	team.Description = request.Description
	team.Avatar = request.Avatar
	team.UpdatedAt = time.Now()
	team.UpdatedBy = authenticatedUser.Email
	if err = service.Repository.Update(id, team); err != nil {
		return err
	}
	return nil
}

func (service *TeamServiceImpl) UpdateStatus(ctx context.Context, id uint, request model.UpdateTeamStatusRequest) error {
	authenticatedUser := ctx.Value("user").(um.User)
	team, err := service.Repository.FindOne(id)
	if err != nil {
		return err
	}

	team.IsActive = request.IsActive
	team.UpdatedBy = authenticatedUser.Email
	team.UpdatedAt = time.Now()
	if err = service.Repository.Update(id, team); err != nil {
		return err
	}
	return nil
}

func (service *TeamServiceImpl) Delete(ctx context.Context, id uint) error {
	authenticatedUser := ctx.Value("user").(um.User)
	if _, err := service.Repository.FindOne(id); err != nil {
		return err
	}

	if err := service.Repository.Delete(id, authenticatedUser.Email); err != nil {
		return err
	}
	return nil
}

func (service *TeamServiceImpl) GetAll(
	filter model.FilterTeam,
	pg *utils.PaginateQueryOffset,
) (response model.GetAllTeamResponse, err error) {
	response.Teams, response.TotalItem, response.TotalPage, err = service.Repository.FindAll(filter, pg)
	if err != nil {
		return
	}
	return
}

func (service *TeamServiceImpl) GetDetail(id uint) (team model.TeamDetail, err error) {
	if team, err = service.Repository.FindDetail(id); err != nil {
		return
	}
	return
}

func (service *TeamServiceImpl) GetListByEventID(
	ctx context.Context,
	filter model.FilterTeam,
	pg *utils.PaginateQueryOffset,
) (response model.GetListTeamByEventIDResponse, err error) {
	authenticatedUser := ctx.Value("user").(um.User)

	if _, err = service.EventRepo.FindOne(filter.EventID); err != nil {
		return
	}

	_, err = service.EventParticipantRepo.FindOneByEventIDAndParticipantID(filter.EventID, authenticatedUser.Participant.ID)
	if err != nil && err != e.ErrDataNotFound {
		return
	} else if err != nil && err == e.ErrDataNotFound {
		err = e.ErrForbidden
		return
	}

	filter.IsForParticipant = true
	filter.Status = "active"
	filter.TeamRequestParticipantID = authenticatedUser.Participant.ID
	teams, totalItem, totalPage, err := service.Repository.FindAll(filter, pg)
	if err != nil {
		return
	}

	for k := range teams {
		response.Teams = append(response.Teams, model.TeamByEventID{
			ID:          teams[k].ID,
			Code:        teams[k].Code,
			Name:        teams[k].Name,
			NumOfMember: teams[k].NumOfMember,
			Avatar:      teams[k].Avatar,
			IsRequested: teams[k].IsRequested,
		})
	}

	response.TotalItem = totalItem
	response.TotalPage = totalPage
	return
}

func (service *TeamServiceImpl) GetDetail2(ctx context.Context, id uint) (team model.TeamDetail2, err error) {
	authenticatedUser := ctx.Value("user").(um.User)
	latestEvent, err := service.EventRepo.FindLatest()
	if err != nil {
		return
	}

	if team, err = service.Repository.FindDetail2(id, latestEvent.ID, authenticatedUser.Participant.ID, true); err != nil {
		return
	}
	return
}

func (service *TeamServiceImpl) GetMyTeam(ctx context.Context) (team model.TeamDetail2, err error) {
	authenticatedUser := ctx.Value("user").(um.User)
	teamMember, err := service.TeamMemberRepo.FindByParticipantID(authenticatedUser.Participant.ID)
	if err != nil {
		return
	}

	latestEvent, err := service.EventRepo.FindLatest()
	if err != nil {
		return
	}

	if team, err = service.Repository.FindDetail2(teamMember.TeamID, latestEvent.ID, 0, false); err != nil {
		return
	}
	return
}

func (service *TeamServiceImpl) GetMembers(ctx context.Context, id uint) (members []model.TeamMemberList, err error) {
	authenticatedUser := ctx.Value("user").(um.User)

	if authenticatedUser.UserRole.Name == constants.UserParticipant {
		_, err = service.TeamMemberRepo.FindByParticipantIDAndTeamID(authenticatedUser.Participant.ID, id)
		if err != nil && err != e.ErrDataNotFound {
			return
		} else if err != nil && err == e.ErrDataNotFound {
			err = e.ErrForbidden
			return
		}
	}

	if members, err = service.TeamMemberRepo.FindManyByTeamID(id); err != nil {
		return
	}
	return
}

func (service *TeamServiceImpl) GetInvitations(ctx context.Context, id uint) (invitations []model.TeamInvitationList, err error) {
	authenticatedUser := ctx.Value("user").(um.User)
	_, err = service.TeamMemberRepo.FindByParticipantIDAndTeamID(authenticatedUser.Participant.ID, id)
	if err != nil && err != e.ErrDataNotFound {
		return
	} else if err != nil && err == e.ErrDataNotFound {
		err = e.ErrForbidden
		return
	}

	latestEvent, err := service.EventRepo.FindLatest()
	if err != nil {
		return
	}

	if invitations, err = service.TeamInvitationRepo.FindManyByTeamIDAndEventID(id, latestEvent.ID); err != nil {
		return
	}
	return
}

func (service *TeamServiceImpl) GetRequests(ctx context.Context, id uint) (requests []model.TeamRequestList, err error) {
	authenticatedUser := ctx.Value("user").(um.User)
	_, err = service.TeamMemberRepo.FindByParticipantIDAndTeamID(authenticatedUser.Participant.ID, id)
	if err != nil && err != e.ErrDataNotFound {
		return
	} else if err != nil && err == e.ErrDataNotFound {
		err = e.ErrForbidden
		return
	}

	latestEvent, err := service.EventRepo.FindLatest()
	if err != nil {
		return
	}

	if requests, err = service.TeamRequestRepo.FindManyByTeamIDAndEventID(id, latestEvent.ID); err != nil {
		return
	}
	return
}
