package service

import (
	"be-sagara-hackathon/src/modules/master-data/speciality/model"
	"be-sagara-hackathon/src/modules/master-data/speciality/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"context"
	"time"
)

type SpecialityService interface {
	CreateSpeciality(ctx context.Context, req model.SpecialityRequest) (err error)
	UpdateSpeciality(ctx context.Context, req model.UpdateSpecialityRequest, id uint) (err error)
	GetListSpeciality(
		filter model.FilterSpeciality,
		pg *utils.PaginateQueryOffset,
	) (response model.ListSpecialityResponse, err error)
	GetDetailSpeciality(id uint) (skill model.Speciality, err error)
}

type SpecialityServiceImpl struct {
	Repository repository.SpecialityRepository
}

func NewSpecialityService(repository repository.SpecialityRepository) SpecialityService {
	return &SpecialityServiceImpl{Repository: repository}
}

func (service *SpecialityServiceImpl) CreateSpeciality(ctx context.Context, req model.SpecialityRequest) (err error) {
	if err = service.Repository.Save(model.Speciality{
		BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
		Name:       req.Name,
	}); err != nil {
		return
	}
	return
}

func (service *SpecialityServiceImpl) UpdateSpeciality(ctx context.Context, req model.UpdateSpecialityRequest, id uint) (err error) {
	speciality, err := service.Repository.FindOne(id)
	if err != nil {
		return
	}

	speciality.Name = req.Name
	speciality.IsActive = req.IsActive
	speciality.UpdatedAt = time.Now()
	speciality.UpdatedBy = ctx.Value("user").(um.User).Email
	if err = service.Repository.Update(speciality, id); err != nil {
		return
	}
	return
}

func (service *SpecialityServiceImpl) GetListSpeciality(filter model.FilterSpeciality, pg *utils.PaginateQueryOffset) (response model.ListSpecialityResponse, err error) {
	specialities, totalData, totalPage, err := service.Repository.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.SpecialityLite
	for _, v := range specialities {
		responseData = append(responseData, model.SpecialityLite{
			ID:       v.ID,
			Name:     v.Name,
			IsActive: v.IsActive,
		})
	}
	response.Specialities = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *SpecialityServiceImpl) GetDetailSpeciality(id uint) (skill model.Speciality, err error) {
	if skill, err = service.Repository.FindOne(id); err != nil {
		return
	}
	return
}
