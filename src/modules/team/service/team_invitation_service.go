package service

import (
	evr "be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/team/model"
	"be-sagara-hackathon/src/modules/team/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	ur "be-sagara-hackathon/src/modules/user/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/constants"
	"be-sagara-hackathon/src/utils/email"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
	"context"
	"fmt"
	"os"
	"time"
)

type TeamInvitationService interface {
	Create(ctx context.Context, request model.CreateInvitationRequest) error
	Update(ctx context.Context, id uint, request model.UpdateInvitationRequest) error
	UpdateStatus(ctx context.Context, code string, request model.UpdateStatusInvitationRequest) error
	Delete(ctx context.Context, id uint) error
	GetList(
		ctx context.Context,
		filter model.FilterInvitation,
		pg *utils.PaginateQueryOffset,
	) (response model.GetListInvitationResponse, err error)
	GetDetail(ctx context.Context, id uint) (invitation model.InvitationDetail, err error)
	GetDetail2(ctx context.Context, id uint) (invitation model.TeamInvitationDetail, err error)
}

type TeamInvitationServiceImpl struct {
	Repository      repository.TeamInvitationRepository
	TeamRepo        repository.TeamRepository
	TeamMemberRepo  repository.TeamMemberRepository
	TeamRequestRepo repository.TeamRequestRepository
	ParticipantRepo ur.ParticipantRepository
	EventRepo       evr.EventRepository
}

func NewTeamInvitationService(
	repository repository.TeamInvitationRepository,
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	teamRequestRepo repository.TeamRequestRepository,
	participantRepo ur.ParticipantRepository,
	eventRepo evr.EventRepository,
) TeamInvitationService {
	return &TeamInvitationServiceImpl{
		Repository:      repository,
		TeamRepo:        teamRepo,
		TeamMemberRepo:  teamMemberRepo,
		TeamRequestRepo: teamRequestRepo,
		ParticipantRepo: participantRepo,
		EventRepo:       eventRepo,
	}
}

func (service *TeamInvitationServiceImpl) Create(ctx context.Context, request model.CreateInvitationRequest) error {
	authenticatedUser := ctx.Value("user").(um.User)
	event, err := service.EventRepo.FindOne(request.EventID)
	if err != nil {
		return err
	}
	if event.Status != constants.EventRunning {
		return e.ErrEventNotRunning
	}

	team, err := service.TeamRepo.FindByIDAndEventID(request.TeamID, request.EventID)
	if err != nil {
		return err
	}
	if team.ParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	if team.NumOfMember >= event.TeamMaxMember {
		return e.ErrTeamIsFull
	}

	invitedParticipant, err := service.ParticipantRepo.FindDetail(request.ToParticipantID)
	if err != nil {
		return err
	}

	if !invitedParticipant.IsRegistered {
		return e.ErrRegistrationNotCompleted
	}

	if invitedParticipant.PaymentStatus != constants.InvoicePaid {
		return e.ErrPaymentNotPaid
	}

	member, err := service.TeamMemberRepo.FindByParticipantID(request.ToParticipantID)
	if err != nil && err != e.ErrDataNotFound {
		return err
	}
	if member.ID != 0 && member.Team.DeletedAt == nil && member.Team.IsActive {
		err = e.ErrHasTeam
		return err
	}

	invs, err := service.Repository.FindByEventIDAndTeamIDAndParticipantID(request.EventID, request.TeamID, request.ToParticipantID)
	if err != nil {
		return err
	}

	for _, v := range invs {
		if v.ID != 0 &&
			v.Status == constants.InvitationOrRequestStatusSent &&
			v.ProceedAt == nil {
			return e.ErrParticipantHasBeenInvited
		}
	}

	teamReqs, err := service.TeamRequestRepo.FindByEventIDAndTeamIDAndParticipantID(event.ID, team.ID, request.ToParticipantID)
	if err != nil {
		return err
	}

	for _, v := range teamReqs {
		if v.ID != 0 &&
			v.Status == constants.InvitationOrRequestStatusSent &&
			v.ProceedAt == nil && v.ProceedBy == nil {
			return e.ErrParticipantRequestedToJoinTeam
		}
	}

	//save invitation
	invitationCode := utils.GenerateUuid()
	if err = service.Repository.Create(model.TeamInvitation{
		BaseEntity:      builder.BuildBaseEntity(ctx, true, nil),
		Code:            invitationCode,
		EventID:         request.EventID,
		TeamID:          request.TeamID,
		ToParticipantID: request.ToParticipantID,
		Status:          constants.InvitationOrRequestStatusSent,
		Note:            request.Note,
	}); err != nil {
		return err
	}

	//TODO:
	// - send email team invitation [OK]
	// - use goroutine/task queue
	templateData := email.TeamInvitationTemplateData{
		Title:       constants.EmailSubjectTeamInvitation,
		InvitedName: invitedParticipant.User.Name,
		SenderName:  authenticatedUser.Name,
		DetailLink:  fmt.Sprintf("%s%s", os.Getenv("TEAM_INVITATION_REDIRECT_URL"), invitationCode),
		AcceptLink:  fmt.Sprintf("%s%s", os.Getenv("ACCEPT_TEAM_INVITATION_REDIRECT_URL"), invitationCode),
	}

	r := email.NewRequest([]string{invitedParticipant.User.Email}, constants.EmailSubjectTeamInvitation, "")
	if err = r.ParseTemplate("./src/utils/email/template_email_team_invitation.html", templateData); err == nil {
		ok, err := r.SendEmail()
		if !ok || err != nil {
			return e.ErrFailedSendEmail
		}
	} else {
		return e.ErrFailedParseEmailTemplate
	}

	return nil
}

func (service *TeamInvitationServiceImpl) Update(ctx context.Context, id uint, request model.UpdateInvitationRequest) error {
	authenticatedUser := ctx.Value("user").(um.User)
	invitation, err := service.Repository.FindOne(id)
	if err != nil {
		return err
	}

	if invitation.Status != constants.InvitationOrRequestStatusSent {
		return e.ErrInvitationHasBeenProceed
	}

	team, err := service.TeamRepo.FindByIDAndEventID(invitation.TeamID, invitation.EventID)
	if err != nil {
		return err
	}

	if team.ParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	if err = service.Repository.Update(id, model.TeamInvitation{
		BaseEntity: common.BaseEntity{
			UpdatedAt: time.Now(),
			UpdatedBy: authenticatedUser.Email,
		},
		Note:      request.Note,
		Status:    invitation.Status,
		ProceedAt: invitation.ProceedAt,
	}, nil); err != nil {
		return err
	}

	return nil
}

func (service *TeamInvitationServiceImpl) UpdateStatus(ctx context.Context, code string, request model.UpdateStatusInvitationRequest) error {
	authenticatedUser := ctx.Value("user").(um.User)
	invitation, err := service.Repository.FindByCode(code)
	if err != nil {
		return err
	}

	if invitation.Status != constants.InvitationOrRequestStatusSent {
		return e.ErrInvitationHasBeenProceed
	}

	team, err := service.TeamRepo.FindByIDAndEventID(invitation.TeamID, invitation.EventID)
	if err != nil {
		return err
	}

	if invitation.ToParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	member, err := service.TeamMemberRepo.FindByParticipantID(authenticatedUser.Participant.ID)
	if err != nil && err != e.ErrDataNotFound {
		return err
	}
	if member.ID != 0 && member.Team.DeletedAt == nil && member.Team.IsActive {
		err = e.ErrHasTeam
		return err
	}

	var newMember *model.TeamMember
	if request.Status == constants.InvitationOrRequestStatusAccepted {
		event, err2 := service.EventRepo.FindOne(invitation.EventID)
		if err2 != nil {
			return err2
		}

		if team.NumOfMember >= event.TeamMaxMember {
			return e.ErrTeamIsFull
		}

		newMember = &model.TeamMember{
			BaseEntity:    builder.BuildBaseEntity(ctx, true, nil),
			TeamID:        invitation.TeamID,
			ParticipantID: authenticatedUser.Participant.ID,
			JoinedAt:      time.Now(),
		}
	}

	updateInvitation := model.TeamInvitation{
		BaseEntity: common.BaseEntity{
			UpdatedAt: time.Now(),
			UpdatedBy: authenticatedUser.Email,
		},
		Status:    request.Status,
		ProceedAt: helper.ReferTime(time.Now()),
		Note:      invitation.Note,
	}
	if err = service.Repository.Update(invitation.ID, updateInvitation, newMember); err != nil {
		return err
	}

	return nil
}

func (service *TeamInvitationServiceImpl) Delete(ctx context.Context, id uint) error {
	authenticatedUser := ctx.Value("user").(um.User)
	invitation, err := service.Repository.FindOne(id)
	if err != nil {
		return err
	}

	if invitation.Status != constants.InvitationOrRequestStatusSent {
		return e.ErrInvitationHasBeenProceed
	}

	team, err := service.TeamRepo.FindByIDAndEventID(invitation.TeamID, invitation.EventID)
	if err != nil {
		return err
	}

	if team.ParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	if err = service.Repository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (service *TeamInvitationServiceImpl) GetList(
	ctx context.Context,
	filter model.FilterInvitation,
	pg *utils.PaginateQueryOffset,
) (response model.GetListInvitationResponse, err error) {
	authenticatedUser := ctx.Value("user").(um.User)
	filter.ParticipantID = authenticatedUser.Participant.ID
	response.Invitations, response.TotalItem, response.TotalPage, err = service.Repository.FindAll(filter, pg)
	if err != nil {
		return
	}
	return
}

func (service *TeamInvitationServiceImpl) GetDetail(ctx context.Context, id uint) (invitation model.InvitationDetail, err error) {
	if invitation, err = service.Repository.FindDetail(id); err != nil {
		return
	}

	authenticatedUser := ctx.Value("user").(um.User)
	if authenticatedUser.Participant.ID != invitation.ToParticipantID {
		err = e.ErrForbidden
		return
	}

	return
}

func (service *TeamInvitationServiceImpl) GetDetail2(ctx context.Context, id uint) (invitation model.TeamInvitationDetail, err error) {
	if invitation, err = service.Repository.FindDetail2(id); err != nil {
		return
	}

	authenticatedUser := ctx.Value("user").(um.User)
	_, err = service.TeamMemberRepo.FindByParticipantIDAndTeamID(authenticatedUser.Participant.ID, invitation.TeamID)
	if err != nil && err != e.ErrDataNotFound {
		return
	} else if err != nil && err == e.ErrDataNotFound {
		err = e.ErrForbidden
		return
	}
	return
}
