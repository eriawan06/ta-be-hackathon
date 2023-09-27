package user

import (
	"be-sagara-hackathon/src/modules/user/controller"
	"be-sagara-hackathon/src/modules/user/repository"
	"be-sagara-hackathon/src/modules/user/service"

	"gorm.io/gorm"
)

var (
	userRoleRepository    repository.UserRoleRepository
	userRepository        repository.UserRepository
	userService           service.UserService
	userController        controller.UserController
	participantRepository repository.ParticipantRepository
	participantService    service.ParticipantService
	participantController controller.ParticipantController
	mentorService         service.MentorService
	mentorController      controller.MentorController
	judgeService          service.JudgeService
	judgeController       controller.JudgeController
)

type Module interface {
	InitModule()
}

type ModuleImpl struct {
	DB *gorm.DB
}

func New(database *gorm.DB) Module {
	return &ModuleImpl{DB: database}
}

func (module ModuleImpl) InitModule() {
	userRoleRepository = repository.NewUserRoleRepository(module.DB)
	userRepository = repository.NewUserRepository(module.DB)
	userService = service.NewUserService(userRepository)
	userController = controller.NewUserController(userService)
	participantRepository = repository.NewParticipantRepository(module.DB)
	participantService = service.NewParticipantService(participantRepository, userRepository, userRoleRepository)
	participantController = controller.NewParticipantController(participantService)
	mentorService = service.NewMentorService(userRepository, userRoleRepository)
	mentorController = controller.NewMentorController(mentorService)
	judgeService = service.NewJudgeService(userRepository, userRoleRepository)
	judgeController = controller.NewJudgeController(judgeService)
}

func GetUserController() controller.UserController {
	return userController
}

func GetUserRepository() repository.UserRepository {
	return userRepository
}

func GetParticipantController() controller.ParticipantController {
	return participantController
}

func GetMentorController() controller.MentorController {
	return mentorController
}

func GetJudgeController() controller.JudgeController {
	return judgeController
}
