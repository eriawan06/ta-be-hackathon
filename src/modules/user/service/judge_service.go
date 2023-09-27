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

type JudgeService interface {
	Create(ctx context.Context, request model.CreateJudgeRequest) error
	Update(ctx context.Context, request model.UpdateJudgeRequest, id uint) error
	Delete(ctx context.Context, id uint) error
	GetList(
		filter model.FilterUser,
		pg *utils.PaginateQueryOffset,
	) (response model.ListJudgeResponse, err error)
	GetDetail(id uint) (judge model.User, err error)
}

type JudgeServiceImpl struct {
	Repository     repository.UserRepository
	RoleRepository repository.UserRoleRepository
}

func NewJudgeService(
	repository repository.UserRepository,
	roleRepository repository.UserRoleRepository,
) JudgeService {
	return &JudgeServiceImpl{
		Repository:     repository,
		RoleRepository: roleRepository,
	}
}

func (service *JudgeServiceImpl) Create(ctx context.Context, request model.CreateJudgeRequest) error {
	//Get Role Judge
	role, err := service.RoleRepository.FindByName(constants.UserJudge)
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

func (service *JudgeServiceImpl) Update(ctx context.Context, request model.UpdateJudgeRequest, id uint) error {
	judge, err := service.Repository.FindByID(id)
	if err != nil {
		return err
	}

	judge.Name = request.Name
	judge.Email = request.Email
	judge.PhoneNumber = &request.PhoneNumber
	judge.Avatar = &request.Avatar
	judge.IsActive = request.IsActive
	judge.OccupationID = &request.OccupationID
	judge.Institution = &request.Institution
	judge.UpdatedAt = time.Now()
	judge.UpdatedBy = ctx.Value("user").(model.User).Email
	if err = service.Repository.Update(id, judge); err != nil {
		return err
	}
	return nil
}

func (service *JudgeServiceImpl) Delete(ctx context.Context, id uint) error {
	authenticatedUser := ctx.Value("user").(model.User)
	if _, err := service.Repository.FindByID(id); err != nil {
		return err
	}

	if err := service.Repository.Delete(id, authenticatedUser.Email); err != nil {
		return err
	}
	return nil
}

func (service *JudgeServiceImpl) GetList(
	filter model.FilterUser,
	pg *utils.PaginateQueryOffset,
) (response model.ListJudgeResponse, err error) {
	//Get Role Judge
	role, err := service.RoleRepository.FindByName(constants.UserJudge)
	if err != nil {
		return
	}

	filter.RoleID = int(role.ID)
	judges, totalData, totalPage, err := service.Repository.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.JudgeLite
	for _, v := range judges {
		var judge model.JudgeLite
		judge.ID = v.ID
		judge.Name = v.Name
		judge.Email = v.Email
		judge.PhoneNumber = v.PhoneNumber
		judge.IsActive = v.IsActive
		responseData = append(responseData, judge)
	}
	response.Judges = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *JudgeServiceImpl) GetDetail(id uint) (judge model.User, err error) {
	if judge, err = service.Repository.FindByID(id); err != nil {
		return
	}
	return
}
