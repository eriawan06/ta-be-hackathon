package service

import (
	"be-sagara-hackathon/src/modules/auth/model"
	um "be-sagara-hackathon/src/modules/user/model"
	ur "be-sagara-hackathon/src/modules/user/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
)

type AdminAuthService interface {
	Login(request model.LoginRequest) (response model.AuthResponse, err error)
}

type AdminAuthServiceImpl struct {
	UserRepo ur.UserRepository
}

func NewAdminAuthService(userRepo ur.UserRepository) AdminAuthService {
	return &AdminAuthServiceImpl{UserRepo: userRepo}
}

func (service *AdminAuthServiceImpl) Login(request model.LoginRequest) (response model.AuthResponse, err error) {
	user, err := service.UserRepo.FindByEmail(request.Email)
	if err != nil {
		if err == e.ErrEmailNotRegistered {
			err = e.ErrWrongLoginCredential
			return
		}
		return
	}

	//check role
	allowedRole := []string{
		constants.UserSuperadmin,
		constants.UserAdmin,
		constants.UserHR,
		constants.UserMentor,
		constants.UserJudge,
	}
	if !helper.StringInSlice(user.UserRole.Name, allowedRole) {
		err = e.ErrForbidden
		return
	}

	isPasswordValid := utils.CheckPasswordHash(request.Password, helper.DereferString(user.Password))
	if !isPasswordValid {
		err = e.ErrWrongLoginCredential
		return
	}

	// Generate Token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return
	}

	response = model.AuthResponse{
		Token: token,
		User: um.UserResponse{
			Id:       user.ID,
			FullName: user.Name,
			Email:    user.Email,
			RoleId:   user.UserRoleID,
			RoleName: user.UserRole.Name,
		},
	}
	return
}
