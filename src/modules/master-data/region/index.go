package region

import (
	"be-sagara-hackathon/src/modules/master-data/region/controller"
	"be-sagara-hackathon/src/modules/master-data/region/repository"
	"be-sagara-hackathon/src/modules/master-data/region/service"
	"gorm.io/gorm"
)

var (
	regionRepository repository.RegionRepository
	regionService    service.RegionService
	regionController controller.RegionController
)

type RegionModule interface {
	InitModule()
}

type RegionModuleImpl struct {
	DB *gorm.DB
}

func New(database *gorm.DB) RegionModule {
	return &RegionModuleImpl{DB: database}
}

func (module *RegionModuleImpl) InitModule() {
	regionRepository = repository.NewRegionRepository(module.DB)
	regionService = service.NewRegionService(regionRepository)
	regionController = controller.NewRegionController(regionService)
}

func GetController() controller.RegionController {
	return regionController
}
