package service

import (
	"be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/modules/user/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
	"context"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, request model.CreateUserRequest) error
	UpdateUser(ctx context.Context, request model.UpdateUserRequest, id uint) error
	DeleteUser(ctx context.Context, id uint) error
	GetList(
		ctx context.Context,
		filter model.FilterUser,
		pg *utils.PaginateQueryOffset,
	) (response model.ListUserResponse, err error)
	GetDetail(userID uint) (user model.User, err error)
	GetUserProfile(ctx context.Context) (profile model.UserProfile, err error)
	ChangePassword(ctx context.Context, request model.ChangePasswordRequest) (err error)
}

type UserServiceImpl struct {
	Repository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		Repository: userRepository,
	}
}

func (service *UserServiceImpl) CreateUser(ctx context.Context, request model.CreateUserRequest) error {
	hashed, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}

	if err = service.Repository.Save(model.User{
		BaseEntity:  builder.BuildBaseEntity(ctx, true, nil),
		UserRoleID:  request.RoleID,
		Name:        request.Name,
		Email:       request.Email,
		PhoneNumber: &request.PhoneNumber,
		Password:    &hashed,
		AuthType:    constants.AuthTypeRegular,
		IsActive:    true,
	}); err != nil {
		return err
	}
	return nil
}

func (service *UserServiceImpl) UpdateUser(ctx context.Context, request model.UpdateUserRequest, id uint) error {
	user, err := service.Repository.FindByID(id)
	if err != nil {
		return err
	}

	user.Name = request.Name
	user.Email = request.Email
	user.PhoneNumber = &request.PhoneNumber
	user.IsActive = request.IsActive
	user.UpdatedAt = time.Now()
	user.UpdatedBy = ctx.Value("user").(model.User).Email
	if err = service.Repository.Update(id, user); err != nil {
		return err
	}
	return nil
}

func (service *UserServiceImpl) DeleteUser(ctx context.Context, id uint) error {
	if _, err := service.Repository.FindByID(id); err != nil {
		return err
	}

	deleteBy := ctx.Value("user").(model.User).Email
	if err := service.Repository.Delete(id, deleteBy); err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) GetList(
	ctx context.Context,
	filter model.FilterUser,
	pg *utils.PaginateQueryOffset,
) (response model.ListUserResponse, err error) {
	authenticatedUser := ctx.Value("user").(model.User)
	filter.ExceptUserID = authenticatedUser.ID
	users, totalData, totalPage, err := service.Repository.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.UserLite
	for _, v := range users {
		responseData = append(responseData, model.UserLite{
			ID:          v.ID,
			Name:        v.Name,
			Email:       v.Email,
			PhoneNumber: v.PhoneNumber,
			IsActive:    v.IsActive,
		})
	}
	response.Users = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *UserServiceImpl) GetDetail(userID uint) (user model.User, err error) {
	if user, err = service.Repository.FindByID(userID); err != nil {
		return
	}
	return
}

func (service *UserServiceImpl) GetUserProfile(ctx context.Context) (profile model.UserProfile, err error) {
	authenticatedUser := ctx.Value("user").(model.User)
	if profile, err = service.Repository.FindUserProfileByEmail(authenticatedUser.Email); err != nil {
		return
	}
	return
}

func (service *UserServiceImpl) ChangePassword(ctx context.Context, request model.ChangePasswordRequest) (err error) {
	authenticatedUser := ctx.Value("user").(model.User)
	user, err := service.Repository.FindByEmail(authenticatedUser.Email)
	if err != nil {
		return
	}

	if request.NewPassword != request.ConfirmNewPassword {
		return e.ErrConfirmPasswordNotSame
	}

	if user.AuthType != constants.AuthTypeRegular {
		return e.ErrWrongAuthMethod
	}

	// validate old password
	isPasswordValid := utils.CheckPasswordHash(request.OldPassword, helper.DereferString(user.Password))
	if !isPasswordValid {
		return e.ErrWrongOldPassword
	}

	hashed, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return
	}

	if err = service.Repository.UpdatePassword(user.ID, hashed, ""); err != nil {
		return
	}
	return
}
