package project

import (
	eve "be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/project/controller"
	"be-sagara-hackathon/src/modules/project/repository"
	"be-sagara-hackathon/src/modules/project/service"
	tm "be-sagara-hackathon/src/modules/team/repository"

	"gorm.io/gorm"
)

var (
	projectRepository repository.ProjectRepository
	projectService    service.ProjectService
	projectController controller.ProjectController

	projectAssessmentRepository repository.ProjectAssessmentRepository
	projectAssessmentService    service.ProjectAssessmentService
	projectAssessmentController controller.ProjectAssessmentController
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
	teamRepository := tm.NewTeamRepository(module.DB)
	teamMemberRepository := tm.NewTeamMemberRepository(module.DB)
	eventRepository := eve.NewEventRepository(module.DB)
	eventJudgeRepository := eve.NewEventJudgeRepository(module.DB)
	criteriaRepository := eve.NewEventAssessmentCriteriaRepository(module.DB)

	projectRepository = repository.NewProjectRepository(module.DB)
	projectService = service.NewProjectService(
		projectRepository,
		teamRepository,
		teamMemberRepository,
		eventRepository,
	)
	projectController = controller.NewProjectController(projectService)

	projectAssessmentRepository = repository.NewProjectAssessmentRepository(module.DB)
	projectAssessmentService = service.NewProjectAssessmentService(
		projectAssessmentRepository,
		projectRepository,
		criteriaRepository,
		eventJudgeRepository,
	)
	projectAssessmentController = controller.NewProjectAssessmentController(projectAssessmentService)
}

func GetProjectController() controller.ProjectController {
	return projectController
}

func GetProjectAssessmentController() controller.ProjectAssessmentController {
	return projectAssessmentController
}
