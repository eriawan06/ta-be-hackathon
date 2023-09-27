package service

import (
	"be-sagara-hackathon/src/modules/master-data/occupation/model"
	"be-sagara-hackathon/src/modules/master-data/occupation/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"context"
	"time"
)

type OccupationService interface {
	CreateOccupation(ctx context.Context, req model.OccupationRequest) (err error)
	UpdateOccupation(ctx context.Context, req model.UpdateOccupationRequest, id uint) (err error)
	GetListOccupation(
		filter model.FilterOccupation,
		pg *utils.PaginateQueryOffset,
	) (response model.ListOccupationResponse, err error)
	GetDetailOccupation(id uint) (occupation model.Occupation, err error)
}

type OccupationServiceImpl struct {
	Repository repository.OccupationRepository
}

func NewOccupationService(repository repository.OccupationRepository) OccupationService {
	return &OccupationServiceImpl{Repository: repository}
}

func (service *OccupationServiceImpl) CreateOccupation(ctx context.Context, req model.OccupationRequest) (err error) {
	if err = service.Repository.Save(model.Occupation{
		BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
		Name:       req.Name,
	}); err != nil {
		return
	}
	return
}

func (service *OccupationServiceImpl) UpdateOccupation(ctx context.Context, req model.UpdateOccupationRequest, id uint) (err error) {
	occupation, err := service.Repository.FindOne(id)
	if err != nil {
		return
	}

	occupation.Name = req.Name
	occupation.IsActive = req.IsActive
	occupation.UpdatedAt = time.Now()
	occupation.UpdatedBy = ctx.Value("user").(um.User).Email
	if err = service.Repository.Update(occupation, id); err != nil {
		return
	}
	return
}

func (service *OccupationServiceImpl) GetListOccupation(filter model.FilterOccupation, pg *utils.PaginateQueryOffset) (response model.ListOccupationResponse, err error) {
	occupations, totalData, totalPage, err := service.Repository.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.OccupationLite
	for _, v := range occupations {
		responseData = append(responseData, model.OccupationLite{
			ID:       v.ID,
			Name:     v.Name,
			IsActive: v.IsActive,
		})
	}
	response.Occupations = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *OccupationServiceImpl) GetDetailOccupation(id uint) (occupation model.Occupation, err error) {
	if occupation, err = service.Repository.FindOne(id); err != nil {
		return
	}
	return
}
