package skill

import (
	"be-sagara-hackathon/src/modules/master-data/skill/controller"
	"be-sagara-hackathon/src/modules/master-data/skill/repository"
	"be-sagara-hackathon/src/modules/master-data/skill/service"
	"gorm.io/gorm"
)

var (
	skillRepository repository.SkillRepository
	skillService    service.SkillService
	skillController controller.SkillController
)

type SkillModule interface {
	InitModule()
}

type SkillModuleImpl struct {
	DB *gorm.DB
}

func New(database *gorm.DB) SkillModule {
	return &SkillModuleImpl{DB: database}
}

func (module *SkillModuleImpl) InitModule() {
	skillRepository = repository.NewSkillRepository(module.DB)
	skillService = service.NewSkillService(skillRepository)
	skillController = controller.NewSkillController(skillService)
}

func GetController() controller.SkillController {
	return skillController
}
