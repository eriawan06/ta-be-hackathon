package service

import (
	"be-sagara-hackathon/src/modules/master-data/technology/model"
	"be-sagara-hackathon/src/modules/master-data/technology/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"context"
	"time"
)

type TechnologyService interface {
	CreateTechnology(ctx context.Context, req model.TechnologyRequest) (err error)
	UpdateTechnology(ctx context.Context, req model.UpdateTechnologyRequest, id uint) (err error)
	GetListTechnology(
		filter model.FilterTechnology,
		pg *utils.PaginateQueryOffset,
	) (response model.ListTechnologyResponse, err error)
	GetDetailTechnology(id uint) (technology model.Technology, err error)
}

type TechnologyServiceImpl struct {
	Repository repository.TechnologyRepository
}

func NewTechnologyService(repository repository.TechnologyRepository) TechnologyService {
	return &TechnologyServiceImpl{Repository: repository}
}

func (service *TechnologyServiceImpl) CreateTechnology(ctx context.Context, req model.TechnologyRequest) (err error) {
	if err = service.Repository.Save(model.Technology{
		BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
		Name:       req.Name,
	}); err != nil {
		return
	}
	return
}

func (service *TechnologyServiceImpl) UpdateTechnology(ctx context.Context, req model.UpdateTechnologyRequest, id uint) (err error) {
	technology, err := service.Repository.FindOne(id)
	if err != nil {
		return
	}

	technology.Name = req.Name
	technology.IsActive = req.IsActive
	technology.UpdatedAt = time.Now()
	technology.UpdatedBy = ctx.Value("user").(um.User).Email
	if err = service.Repository.Update(technology, id); err != nil {
		return
	}
	return
}

func (service *TechnologyServiceImpl) GetListTechnology(
	filter model.FilterTechnology,
	pg *utils.PaginateQueryOffset,
) (response model.ListTechnologyResponse, err error) {
	technologies, totalData, totalPage, err := service.Repository.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.TechnologyLite
	for _, v := range technologies {
		responseData = append(responseData, model.TechnologyLite{
			ID:       v.ID,
			Name:     v.Name,
			IsActive: v.IsActive,
		})
	}
	response.Technologies = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *TechnologyServiceImpl) GetDetailTechnology(id uint) (technology model.Technology, err error) {
	if technology, err = service.Repository.FindOne(id); err != nil {
		return
	}
	return
}
