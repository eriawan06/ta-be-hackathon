package team

import (
	ever "be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/team/controller"
	"be-sagara-hackathon/src/modules/team/repository"
	"be-sagara-hackathon/src/modules/team/service"
	ur "be-sagara-hackathon/src/modules/user/repository"
	"gorm.io/gorm"
)

var (
	teamRepository repository.TeamRepository
	teamService    service.TeamService
	teamController controller.TeamController

	teamInvitationRepository repository.TeamInvitationRepository
	teamInvitationService    service.TeamInvitationService
	teamInvitationController controller.TeamInvitationController

	teamRequestRepository repository.TeamRequestRepository
	teamRequestService    service.TeamRequestService
	teamRequestController controller.TeamRequestController

	teamMemberRepository repository.TeamMemberRepository
	teamMemberService    service.TeamMemberService
	teamMemberController controller.TeamMemberController
)

type Module interface {
	InitModule()
}

type ModuleImpl struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Module {
	return &ModuleImpl{DB: db}
}

func (module ModuleImpl) InitModule() {
	participantRepository := ur.NewParticipantRepository(module.DB)
	eventRepository := ever.NewEventRepository(module.DB)
	eventParticipantRepository := ever.NewEventParticipantRepository(module.DB)
	teamMemberRepository = repository.NewTeamMemberRepository(module.DB)
	teamInvitationRepository = repository.NewTeamInvitationRepository(module.DB)
	teamRequestRepository = repository.NewTeamRequestRepository(module.DB)
	teamRepository = repository.NewTeamRepository(module.DB)

	teamService = service.NewTeamService(
		teamRepository,
		teamMemberRepository,
		teamInvitationRepository,
		teamRequestRepository,
		participantRepository,
		eventRepository,
		eventParticipantRepository,
	)
	teamController = controller.NewTeamController(teamService)

	teamInvitationService = service.NewTeamInvitationService(
		teamInvitationRepository,
		teamRepository,
		teamMemberRepository,
		teamRequestRepository,
		participantRepository,
		eventRepository,
	)
	teamInvitationController = controller.NewTeamInvitationController(teamInvitationService)

	teamRequestService = service.NewTeamRequestService(
		teamRequestRepository,
		teamRepository,
		teamMemberRepository,
		teamInvitationRepository,
		participantRepository,
		eventRepository,
	)
	teamRequestController = controller.NewTeamRequestController(teamRequestService)

	teamMemberService = service.NewTeamMemberService(teamMemberRepository)
	teamMemberController = controller.NewTeamMemberController(teamMemberService)
}

func GetTeamController() controller.TeamController {
	return teamController
}

func GetTeamInvitationController() controller.TeamInvitationController {
	return teamInvitationController
}

func GetTeamMemberController() controller.TeamMemberController {
	return teamMemberController
}

func GetTeamRequestController() controller.TeamRequestController {
	return teamRequestController
}
