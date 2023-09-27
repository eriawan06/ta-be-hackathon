package event

import (
	"be-sagara-hackathon/src/modules/event/controller"
	"be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/event/service"
	scr "be-sagara-hackathon/src/modules/schedule/repository"
	tr "be-sagara-hackathon/src/modules/team/repository"

	"gorm.io/gorm"
)

var (
	eventRepository            repository.EventRepository
	eventParticipantRepository repository.EventParticipantRepository
	eventService               service.EventService

	eventController                   controller.EventController
	eventMentorController             controller.EventMentorController
	eventJudgeController              controller.EventJudgeController
	eventCompanyController            controller.EventCompanyController
	eventTimelineController           controller.EventTimelineController
	eventRuleController               controller.EventRuleController
	eventFaqController                controller.EventFaqController
	eventAssessmentCriteriaController controller.EventAssessmentCriteriaController
)

type EventModule interface {
	InitModule()
}

type EventModuleImpl struct {
	DB *gorm.DB
}

func New(database *gorm.DB) EventModule {
	return &EventModuleImpl{DB: database}
}

func (module *EventModuleImpl) InitModule() {
	teamMemberRepo := tr.NewTeamMemberRepository(module.DB)
	scheduleRepo := scr.NewScheduleRepository(module.DB)

	eventParticipantRepository = repository.NewEventParticipantRepository(module.DB)
	eventRepository = repository.NewEventRepository(module.DB)
	eventService = service.NewEventRepository(eventRepository, eventParticipantRepository, teamMemberRepo, scheduleRepo)
	eventController = controller.NewEventController(eventService)

	eventMentorRepository := repository.NewEventMentorRepository(module.DB)
	eventMentorService := service.NewEventMentorService(eventMentorRepository)
	eventMentorController = controller.NewEventMentorController(eventMentorService)

	eventJudgeRepository := repository.NewEventJudgeRepository(module.DB)
	eventJudgeService := service.NewEventJudgeService(eventJudgeRepository)
	eventJudgeController = controller.NewEventJudgeController(eventJudgeService)

	eventCompanyRepository := repository.NewEventCompanyRepository(module.DB)
	eventCompanyService := service.NewEventCompanyService(eventCompanyRepository, eventRepository)
	eventCompanyController = controller.NewEventCompanyController(eventCompanyService)

	eventTimelineRepository := repository.NewEventTimelineRepository(module.DB)
	eventTimelineService := service.NewEventTimelineService(eventTimelineRepository, eventRepository)
	eventTimelineController = controller.NewEventTimelineController(eventTimelineService)

	eventRuleRepository := repository.NewEventRuleRepository(module.DB)
	eventRuleService := service.NewEventRuleService(eventRuleRepository, eventRepository)
	eventRuleController = controller.NewEventRuleController(eventRuleService)

	eventFaqRepository := repository.NewEventFaqRepository(module.DB)
	eventFaqService := service.NewEventFaqService(eventFaqRepository, eventRepository)
	eventFaqController = controller.NewEventFaqController(eventFaqService)

	eventAssessmentCriteriaRepository := repository.NewEventAssessmentCriteriaRepository(module.DB)
	eventAssessmentCriteriaService := service.NewEventAssessmentCriteriaService(eventAssessmentCriteriaRepository, eventRepository)
	eventAssessmentCriteriaController = controller.NewEventAssessmentCriteriaController(eventAssessmentCriteriaService)
}

func GetController() controller.EventController {
	return eventController
}

func GetEventJudgeController() controller.EventJudgeController {
	return eventJudgeController
}

func GetEventMentorController() controller.EventMentorController {
	return eventMentorController
}

func GetEventCompanyController() controller.EventCompanyController {
	return eventCompanyController
}

func GetEventTimelineController() controller.EventTimelineController {
	return eventTimelineController
}

func GetEventRuleController() controller.EventRuleController {
	return eventRuleController
}

func GetEventFaqController() controller.EventFaqController {
	return eventFaqController
}

func GetEventAssessmentCriteriaController() controller.EventAssessmentCriteriaController {
	return eventAssessmentCriteriaController
}

func GetService() service.EventService {
	return eventService
}
