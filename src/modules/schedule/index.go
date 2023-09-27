package schedule

import (
	evr "be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/schedule/controller"
	"be-sagara-hackathon/src/modules/schedule/repository"
	"be-sagara-hackathon/src/modules/schedule/service"
	tr "be-sagara-hackathon/src/modules/team/repository"
	ur "be-sagara-hackathon/src/modules/user/repository"
	"gorm.io/gorm"
)

var (
	scheduleRepository repository.ScheduleRepository
	scheduleService    service.ScheduleService
	scheduleController controller.ScheduleController
)

type ScheduleModule interface {
	InitModule()
}

type ScheduleModuleImpl struct {
	DB *gorm.DB
}

func New(database *gorm.DB) ScheduleModule {
	return &ScheduleModuleImpl{DB: database}
}

func (module *ScheduleModuleImpl) InitModule() {
	eventRepository := evr.NewEventRepository(module.DB)
	userRepository := ur.NewUserRepository(module.DB)
	teamRepository := tr.NewTeamRepository(module.DB)
	scheduleRepository = repository.NewScheduleRepository(module.DB)
	scheduleService = service.NewScheduleService(scheduleRepository, eventRepository, userRepository, teamRepository)
	scheduleController = controller.NewScheduleController(scheduleService)
}

func GetController() controller.ScheduleController {
	return scheduleController
}
