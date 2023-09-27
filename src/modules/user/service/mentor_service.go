package service

import (
	"be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/modules/user/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/constants"
	"context"
	"time"
)

type MentorService interface {
	Create(ctx context.Context, request model.CreateMentorRequest) error
	Update(ctx context.Context, request model.UpdateMentorRequest, id uint) error
	Delete(ctx context.Context, id uint) error
	GetList(
		filter model.FilterUser,
		pg *utils.PaginateQueryOffset,
	) (response model.ListMentorResponse, err error)
	GetDetail(id uint) (mentor model.User, err error)
}

type MentorServiceImpl struct {
	Repository     repository.UserRepository
	RoleRepository repository.UserRoleRepository
}

func NewMentorService(
	repository repository.UserRepository,
	roleRepository repository.UserRoleRepository,
) MentorService {
	return &MentorServiceImpl{
		Repository:     repository,
		RoleRepository: roleRepository,
	}
}

func (service *MentorServiceImpl) Create(ctx context.Context, request model.CreateMentorRequest) error {
	//Get Role Mentor
	role, err := service.RoleRepository.FindByName(constants.UserMentor)
	if err != nil {
		return err
	}

	hashed, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}

	if err = service.Repository.Save(model.User{
		BaseEntity:   builder.BuildBaseEntity(ctx, true, nil),
		UserRoleID:   role.ID,
		Name:         request.Name,
		Email:        request.Email,
		PhoneNumber:  &request.PhoneNumber,
		Password:     &hashed,
		Avatar:       &request.Avatar,
		AuthType:     constants.AuthTypeRegular,
		IsActive:     true,
		OccupationID: &request.OccupationID,
		Institution:  &request.Institution,
	}); err != nil {
		return err
	}
	return nil
}

func (service *MentorServiceImpl) Update(ctx context.Context, request model.UpdateMentorRequest, id uint) error {
	mentor, err := service.Repository.FindByID(id)
	if err != nil {
		return err
	}

	mentor.Name = request.Name
	mentor.Email = request.Email
	mentor.PhoneNumber = &request.PhoneNumber
	mentor.Avatar = &request.Avatar
	mentor.IsActive = request.IsActive
	mentor.OccupationID = &request.OccupationID
	mentor.Institution = &request.Institution
	mentor.UpdatedAt = time.Now()
	mentor.UpdatedBy = ctx.Value("user").(model.User).Email
	if err = service.Repository.Update(id, mentor); err != nil {
		return err
	}
	return nil
}

func (service *MentorServiceImpl) Delete(ctx context.Context, id uint) error {
	authenticatedUser := ctx.Value("user").(model.User)
	if _, err := service.Repository.FindByID(id); err != nil {
		return err
	}

	if err := service.Repository.Delete(id, authenticatedUser.Email); err != nil {
		return err
	}
	return nil
}

func (service *MentorServiceImpl) GetList(
	filter model.FilterUser,
	pg *utils.PaginateQueryOffset,
) (response model.ListMentorResponse, err error) {
	//Get Role Mentor
	role, err := service.RoleRepository.FindByName(constants.UserMentor)
	if err != nil {
		return
	}

	filter.RoleID = int(role.ID)
	mentors, totalData, totalPage, err := service.Repository.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.MentorLite
	for _, v := range mentors {
		var mentor model.MentorLite
		mentor.ID = v.ID
		mentor.Name = v.Name
		mentor.Email = v.Email
		mentor.PhoneNumber = v.PhoneNumber
		mentor.IsActive = v.IsActive
		responseData = append(responseData, mentor)
	}
	response.Mentors = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *MentorServiceImpl) GetDetail(id uint) (mentor model.User, err error) {
	if mentor, err = service.Repository.FindByID(id); err != nil {
		return
	}
	return
}
