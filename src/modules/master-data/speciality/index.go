package speciality

import (
	"be-sagara-hackathon/src/modules/master-data/speciality/controller"
	"be-sagara-hackathon/src/modules/master-data/speciality/repository"
	"be-sagara-hackathon/src/modules/master-data/speciality/service"
	"gorm.io/gorm"
)

var (
	specialityRepository repository.SpecialityRepository
	specialityService    service.SpecialityService
	specialityController controller.SpecialityController
)

type SpecialityModule interface {
	InitModule()
}

type SpecialityModuleImpl struct {
	DB *gorm.DB
}

func New(database *gorm.DB) SpecialityModule {
	return &SpecialityModuleImpl{DB: database}
}

func (module *SpecialityModuleImpl) InitModule() {
	specialityRepository = repository.NewSpecialityRepository(module.DB)
	specialityService = service.NewSpecialityService(specialityRepository)
	specialityController = controller.NewSpecialityController(specialityService)
}

func GetController() controller.SpecialityController {
	return specialityController
}
