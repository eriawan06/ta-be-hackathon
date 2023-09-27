package service

import (
	eve "be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/project/model"
	"be-sagara-hackathon/src/modules/project/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"context"
)

type ProjectAssessmentService interface {
	CreateBatch(ctx context.Context, projectID uint, request model.CreateBatchProjectAssessmentRequest) error
	GetByProjectID(projectID uint) (assessments []model.GetByProjectIDResponse, err error)
	GetByJudgeAndProjectID(ctx context.Context, projectID uint) (assessments []model.ProjectAssessment, err error)
}

type ProjectAssessmentServiceImpl struct {
	Repository     repository.ProjectAssessmentRepository
	ProjectRepo    repository.ProjectRepository
	CriteriaRepo   eve.EventAssessmentCriteriaRepository
	EventJudgeRepo eve.EventJudgeRepository
}

func NewProjectAssessmentService(
	repo repository.ProjectAssessmentRepository,
	projectRepo repository.ProjectRepository,
	criteriaRepo eve.EventAssessmentCriteriaRepository,
	eventJudgeRepo eve.EventJudgeRepository,
) ProjectAssessmentService {
	return &ProjectAssessmentServiceImpl{
		Repository:     repo,
		ProjectRepo:    projectRepo,
		CriteriaRepo:   criteriaRepo,
		EventJudgeRepo: eventJudgeRepo,
	}
}

func (service *ProjectAssessmentServiceImpl) CreateBatch(ctx context.Context, projectID uint, request model.CreateBatchProjectAssessmentRequest) error {
	authenticatedUser := ctx.Value("user").(um.User)
	project, err := service.ProjectRepo.FindOne(projectID)
	if err != nil {
		return err
	}

	if project.Status != constants.ProjectStatusSubmitted {
		return e.ErrProjectStatusShouldBeSubmitted
	}

	_, err = service.EventJudgeRepo.FindOneByJudgeIDAndEventID(authenticatedUser.ID, project.EventID)
	if err != nil && err != e.ErrDataNotFound {
		return err
	} else if err != nil && err == e.ErrDataNotFound {
		return e.ErrForbidden
	}

	var data []model.ProjectAssessment
	for k, v := range request.Assessments {
		if _, err = service.CriteriaRepo.FindOne(v.CriteriaID); err != nil {
			return err
		}

		data = append(data, model.ProjectAssessment{
			BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
			JudgeID:    authenticatedUser.ID,
			ProjectID:  projectID,
			CriteriaID: request.Assessments[k].CriteriaID,
			Score:      request.Assessments[k].Score,
		})
	}

	if err = service.Repository.CreateBatch(data); err != nil {
		return err
	}
	return nil
}

func (service *ProjectAssessmentServiceImpl) GetByProjectID(projectID uint) (assessments []model.GetByProjectIDResponse, err error) {
	data, err := service.Repository.FindByProjectID(projectID)
	if err != nil {
		return
	}

	isExist := false
	for k, v := range data {
		for k2, v2 := range assessments {
			if v2.JudgeID == v.JudgeID {
				assessments[k2].Assessments = append(assessments[k2].Assessments, data[k])
				isExist = true
			}
		}

		if !isExist {
			assessments = append(assessments, model.GetByProjectIDResponse{
				JudgeID:     data[k].JudgeID,
				JudgeName:   data[k].Judge.Name,
				Assessments: []model.ProjectAssessment{data[k]},
			})
		}
		isExist = false
	}

	return
}

func (service *ProjectAssessmentServiceImpl) GetByJudgeAndProjectID(ctx context.Context, projectID uint) (assessments []model.ProjectAssessment, err error) {
	authenticatedUser := ctx.Value("user").(um.User)
	project, err := service.ProjectRepo.FindOne(projectID)
	if err != nil {
		return
	}

	_, err = service.EventJudgeRepo.FindOneByJudgeIDAndEventID(authenticatedUser.ID, project.EventID)
	if err != nil && err != e.ErrDataNotFound {
		return
	} else if err != nil && err == e.ErrDataNotFound {
		err = e.ErrForbidden
		return
	}

	assessments, err = service.Repository.FindByProjectIDAndJudgeID(projectID, authenticatedUser.ID)
	if err != nil {
		return
	}

	return
}
