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

type TeamRequestService interface {
	Create(ctx context.Context, request model.CreateRequestJoinTeam) error
	Update(ctx context.Context, id uint, request model.UpdateRequestJoinTeam) error
	UpdateStatus(ctx context.Context, code string, request model.UpdateStatusRequestJoinTeam) error
	Delete(ctx context.Context, id uint) error
	GetDetail(ctx context.Context, id uint) (teamReq model.TeamRequest, err error)
	GetDetailFull(ctx context.Context, id uint) (teamReq model.TeamRequestDetail, err error)
}

type TeamRequestServiceImpl struct {
	Repository         repository.TeamRequestRepository
	TeamRepo           repository.TeamRepository
	TeamMemberRepo     repository.TeamMemberRepository
	TeamInvitationRepo repository.TeamInvitationRepository
	ParticipantRepo    ur.ParticipantRepository
	EventRepo          evr.EventRepository
}

func NewTeamRequestService(
	repository repository.TeamRequestRepository,
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	teamInvitationRepo repository.TeamInvitationRepository,
	participantRepo ur.ParticipantRepository,
	eventRepo evr.EventRepository,
) TeamRequestService {
	return &TeamRequestServiceImpl{
		Repository:         repository,
		TeamRepo:           teamRepo,
		TeamMemberRepo:     teamMemberRepo,
		TeamInvitationRepo: teamInvitationRepo,
		ParticipantRepo:    participantRepo,
		EventRepo:          eventRepo,
	}
}

func (service *TeamRequestServiceImpl) Create(ctx context.Context, request model.CreateRequestJoinTeam) error {
	authenticatedUser := ctx.Value("user").(um.User)
	if !authenticatedUser.Participant.IsRegistered {
		return e.ErrRegistrationNotCompleted
	}

	if authenticatedUser.Participant.PaymentStatus != constants.InvoicePaid {
		return e.ErrPaymentNotPaid
	}

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

	if team.ParticipantID == authenticatedUser.Participant.ID {
		return e.ErrCannotRequestToJoinYourTeam
	}

	if team.NumOfMember >= event.TeamMaxMember {
		return e.ErrTeamIsFull
	}

	member, err := service.TeamMemberRepo.FindByParticipantID(authenticatedUser.Participant.ID)
	if err != nil && err != e.ErrDataNotFound {
		return err
	}
	if member.ID != 0 && member.Team.DeletedAt == nil && member.Team.IsActive {
		err = e.ErrHasTeam
		return err
	}

	teamReqs, err := service.Repository.FindByEventIDAndTeamIDAndParticipantID(event.ID, team.ID, authenticatedUser.Participant.ID)
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

	invs, err := service.TeamInvitationRepo.FindByEventIDAndTeamIDAndParticipantID(request.EventID, request.TeamID, authenticatedUser.Participant.ID)
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

	//save request
	requestCode := utils.GenerateUuid()
	if err = service.Repository.Create(model.TeamRequest{
		BaseEntity:    builder.BuildBaseEntity(ctx, true, nil),
		Code:          requestCode,
		EventID:       request.EventID,
		TeamID:        request.TeamID,
		ParticipantID: authenticatedUser.Participant.ID,
		Status:        constants.InvitationOrRequestStatusSent,
		Note:          request.Note,
	}); err != nil {
		return err
	}

	//TODO:
	// - send email team request [OK]
	// - use goroutine/task queue
	templateData := email.TeamRequestTemplateData{
		Title:           constants.EmailSubjectTeamRequest,
		TeamCreatorName: team.ParticipantName,
		SenderName:      authenticatedUser.Name,
		DetailLink:      fmt.Sprintf("%s%s", os.Getenv("TEAM_REQUEST_REDIRECT_URL"), requestCode),
		AcceptLink:      fmt.Sprintf("%s%s", os.Getenv("ACCEPT_TEAM_REQUEST_REDIRECT_URL"), requestCode),
	}

	r := email.NewRequest([]string{team.ParticipantEmail}, constants.EmailSubjectTeamRequest, "")
	if err = r.ParseTemplate("./src/utils/email/template_email_team_request.html", templateData); err == nil {
		ok, err := r.SendEmail()
		if !ok || err != nil {
			return e.ErrFailedSendEmail
		}
	} else {
		return e.ErrFailedParseEmailTemplate
	}

	return nil
}

func (service *TeamRequestServiceImpl) Update(ctx context.Context, id uint, request model.UpdateRequestJoinTeam) error {
	authenticatedUser := ctx.Value("user").(um.User)
	teamReq, err := service.Repository.FindOne(id)
	if err != nil {
		return err
	}

	if teamReq.Status != constants.InvitationOrRequestStatusSent {
		return e.ErrTeamReqHasBeenProceed
	}

	if teamReq.ParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	if err = service.Repository.Update(id, model.TeamRequest{
		BaseEntity: common.BaseEntity{
			UpdatedAt: time.Now(),
			UpdatedBy: authenticatedUser.Email,
		},
		Note:      request.Note,
		Status:    teamReq.Status,
		ProceedAt: teamReq.ProceedAt,
		ProceedBy: teamReq.ProceedBy,
	}, nil); err != nil {
		return err
	}

	return nil
}

func (service *TeamRequestServiceImpl) UpdateStatus(ctx context.Context, code string, request model.UpdateStatusRequestJoinTeam) error {
	authenticatedUser := ctx.Value("user").(um.User)
	teamReq, err := service.Repository.FindByCode(code)
	if err != nil {
		return err
	}

	if teamReq.Status != constants.InvitationOrRequestStatusSent {
		return e.ErrTeamReqHasBeenProceed
	}

	team, err := service.TeamRepo.FindByIDAndEventID(teamReq.TeamID, teamReq.EventID)
	if err != nil {
		return err
	}

	if team.ParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	member, err := service.TeamMemberRepo.FindByParticipantID(teamReq.ParticipantID)
	if err != nil && err != e.ErrDataNotFound {
		return err
	}
	if member.ID != 0 && member.Team.DeletedAt == nil && member.Team.IsActive {
		err = e.ErrHasTeam
		return err
	}

	var newMember *model.TeamMember
	if request.Status == constants.InvitationOrRequestStatusAccepted {
		event, err2 := service.EventRepo.FindOne(teamReq.EventID)
		if err2 != nil {
			return err2
		}

		if team.NumOfMember >= event.TeamMaxMember {
			return e.ErrTeamIsFull
		}

		newMember = &model.TeamMember{
			BaseEntity:    builder.BuildBaseEntity(ctx, true, nil),
			TeamID:        teamReq.TeamID,
			ParticipantID: teamReq.ParticipantID,
			JoinedAt:      time.Now(),
		}
	}

	updateTeamReq := model.TeamRequest{
		BaseEntity: common.BaseEntity{
			UpdatedAt: time.Now(),
			UpdatedBy: authenticatedUser.Email,
		},
		Status:    request.Status,
		ProceedAt: helper.ReferTime(time.Now()),
		ProceedBy: &authenticatedUser.Email,
		Note:      teamReq.Note,
	}
	if err = service.Repository.Update(teamReq.ID, updateTeamReq, newMember); err != nil {
		return err
	}

	return nil
}

func (service *TeamRequestServiceImpl) Delete(ctx context.Context, id uint) error {
	authenticatedUser := ctx.Value("user").(um.User)
	teamReq, err := service.Repository.FindOne(id)
	if err != nil {
		return err
	}

	if teamReq.Status != constants.InvitationOrRequestStatusSent {
		return e.ErrTeamReqHasBeenProceed
	}

	if teamReq.ParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	if err = service.Repository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (service *TeamRequestServiceImpl) GetDetail(ctx context.Context, id uint) (teamReq model.TeamRequest, err error) {
	if teamReq, err = service.Repository.FindOne(id); err != nil {
		return
	}

	authenticatedUser := ctx.Value("user").(um.User)
	if authenticatedUser.Participant.ID != teamReq.ParticipantID {
		err = e.ErrForbidden
		return
	}

	return
}

func (service *TeamRequestServiceImpl) GetDetailFull(ctx context.Context, id uint) (teamReq model.TeamRequestDetail, err error) {
	if teamReq, err = service.Repository.FindDetail(id); err != nil {
		return
	}

	authenticatedUser := ctx.Value("user").(um.User)
	_, err = service.TeamMemberRepo.FindByParticipantIDAndTeamID(authenticatedUser.Participant.ID, teamReq.TeamID)
	if err != nil && err != e.ErrDataNotFound {
		return
	} else if err != nil && err == e.ErrDataNotFound {
		err = e.ErrForbidden
		return
	}
	return
}
