package auth

import (
	"be-sagara-hackathon/src/modules/auth/controller"
	"be-sagara-hackathon/src/modules/auth/repository"
	"be-sagara-hackathon/src/modules/auth/service"
	er "be-sagara-hackathon/src/modules/event/repository"
	ur "be-sagara-hackathon/src/modules/user/repository"

	"gorm.io/gorm"
)

var (
	// userService    = user.GetService()
	authRepository      repository.AuthRepository
	authService         service.AuthService
	authController      controller.AuthController
	adminAuthService    service.AdminAuthService
	adminAuthController controller.AdminAuthController
)

type AuthModule interface {
	InitModule()
}

type AuthModuleImpl struct {
	DB *gorm.DB
}

func New(database *gorm.DB) AuthModule {
	return &AuthModuleImpl{DB: database}
}

func (module *AuthModuleImpl) InitModule() {
	authRepository = repository.NewAuthRepository(module.DB)
	verifCodeRepository := repository.NewVerificationCodeRepository(module.DB)
	userRepository := ur.NewUserRepository(module.DB)
	participantRepository := ur.NewParticipantRepository(module.DB)
	eventRepository := er.NewEventRepository(module.DB)
	userRoleRepository := ur.NewUserRoleRepository(module.DB)
	eventParticipantRepository := er.NewEventParticipantRepository(module.DB)
	authService = service.NewAuthService(
		authRepository,
		verifCodeRepository,
		userRepository,
		participantRepository,
		eventRepository,
		userRoleRepository,
		eventParticipantRepository,
	)
	authController = controller.NewAuthController(authService)
	adminAuthService = service.NewAdminAuthService(userRepository)
	adminAuthController = controller.NewAdminAuthController(adminAuthService)
}

func GetController() controller.AuthController {
	return authController
}

func GetAdminController() controller.AdminAuthController {
	return adminAuthController
}
