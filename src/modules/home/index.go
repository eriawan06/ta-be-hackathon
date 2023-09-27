package home

import (
	er "be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/home/controller"
	"be-sagara-hackathon/src/modules/home/service"
	"gorm.io/gorm"
)

var (
	homeService    service.HomeService
	homeController controller.HomeController
)

type HomeModule interface {
	InitModule()
}

type HomeModuleImpl struct {
	DB *gorm.DB
}

func New(database *gorm.DB) HomeModule {
	return &HomeModuleImpl{DB: database}
}

func (module *HomeModuleImpl) InitModule() {
	eventRepository := er.NewEventRepository(module.DB)
	eventJudgeRepository := er.NewEventJudgeRepository(module.DB)
	eventMentorRepository := er.NewEventMentorRepository(module.DB)
	eventTimelineRepository := er.NewEventTimelineRepository(module.DB)
	eventFaqRepository := er.NewEventFaqRepository(module.DB)
	eventCompanyRepository := er.NewEventCompanyRepository(module.DB)
	homeService = service.NewHomeService(
		eventRepository,
		eventJudgeRepository,
		eventMentorRepository,
		eventTimelineRepository,
		eventFaqRepository,
		eventCompanyRepository,
	)
	homeController = controller.NewHomeController(homeService)
}

func GetController() controller.HomeController {
	return homeController
}
