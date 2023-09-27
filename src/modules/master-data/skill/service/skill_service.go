package service

import (
	"be-sagara-hackathon/src/modules/master-data/skill/model"
	"be-sagara-hackathon/src/modules/master-data/skill/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"context"
	"time"
)

type SkillService interface {
	CreateSkill(ctx context.Context, req model.SkillRequest) (err error)
	UpdateSkill(ctx context.Context, req model.UpdateSkillRequest, id uint) (err error)
	GetListSkill(
		filter model.FilterSkill,
		pg *utils.PaginateQueryOffset,
	) (response model.ListSkillResponse, err error)
	GetDetailSkill(id uint) (skill model.Skill, err error)
}

type SkillServiceImpl struct {
	Repository repository.SkillRepository
}

func NewSkillService(repository repository.SkillRepository) SkillService {
	return &SkillServiceImpl{Repository: repository}
}

func (service *SkillServiceImpl) CreateSkill(ctx context.Context, req model.SkillRequest) (err error) {
	if err = service.Repository.Save(model.Skill{
		BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
		Name:       req.Name,
	}); err != nil {
		return
	}
	return
}

func (service *SkillServiceImpl) UpdateSkill(ctx context.Context, req model.UpdateSkillRequest, id uint) (err error) {
	skill, err := service.Repository.FindOne(id)
	if err != nil {
		return
	}

	skill.Name = req.Name
	skill.IsActive = req.IsActive
	skill.UpdatedAt = time.Now()
	skill.UpdatedBy = ctx.Value("user").(um.User).Email
	if err = service.Repository.Update(skill, id); err != nil {
		return
	}
	return
}

func (service *SkillServiceImpl) GetListSkill(filter model.FilterSkill, pg *utils.PaginateQueryOffset) (response model.ListSkillResponse, err error) {
	skills, totalData, totalPage, err := service.Repository.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.SkillLite
	for _, v := range skills {
		responseData = append(responseData, model.SkillLite{
			ID:       v.ID,
			Name:     v.Name,
			IsActive: v.IsActive,
		})
	}
	response.Skills = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *SkillServiceImpl) GetDetailSkill(id uint) (skill model.Skill, err error) {
	if skill, err = service.Repository.FindOne(id); err != nil {
		return
	}
	return
}
