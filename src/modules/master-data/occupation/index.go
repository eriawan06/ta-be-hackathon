package occupation

import (
	"be-sagara-hackathon/src/modules/master-data/occupation/controller"
	"be-sagara-hackathon/src/modules/master-data/occupation/repository"
	"be-sagara-hackathon/src/modules/master-data/occupation/service"
	"gorm.io/gorm"
)

var (
	occupationRepository repository.OccupationRepository
	occupationService    service.OccupationService
	occupationController controller.OccupationController
)

type OccupationModule interface {
	InitModule()
}

type OccupationModuleImpl struct {
	DB *gorm.DB
}

func New(database *gorm.DB) OccupationModule {
	return &OccupationModuleImpl{DB: database}
}

func (module *OccupationModuleImpl) InitModule() {
	occupationRepository = repository.NewOccupationRepository(module.DB)
	occupationService = service.NewOccupationService(occupationRepository)
	occupationController = controller.NewOccupationController(occupationService)
}

func GetController() controller.OccupationController {
	return occupationController
}
