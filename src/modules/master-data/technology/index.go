package technology

import (
	"be-sagara-hackathon/src/modules/master-data/technology/controller"
	"be-sagara-hackathon/src/modules/master-data/technology/repository"
	"be-sagara-hackathon/src/modules/master-data/technology/service"
	"gorm.io/gorm"
)

var (
	technologyRepository repository.TechnologyRepository
	technologyService    service.TechnologyService
	technologyController controller.TechnologyController
)

type TechnologyModule interface {
	InitModule()
}

type TechnologyModuleImpl struct {
	DB *gorm.DB
}

func New(database *gorm.DB) TechnologyModule {
	return &TechnologyModuleImpl{DB: database}
}

func (module *TechnologyModuleImpl) InitModule() {
	technologyRepository = repository.NewTechnologyRepository(module.DB)
	technologyService = service.NewTechnologyService(technologyRepository)
	technologyController = controller.NewTechnologyController(technologyService)
}

func GetController() controller.TechnologyController {
	return technologyController
}
