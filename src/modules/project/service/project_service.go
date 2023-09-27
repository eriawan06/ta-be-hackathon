package service

import (
	eve "be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/project/model"
	"be-sagara-hackathon/src/modules/project/repository"
	tm "be-sagara-hackathon/src/modules/team/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
	"context"
	"time"
)

type ProjectService interface {
	Create(ctx context.Context, request model.CreateProjectRequest) error
	Update(ctx context.Context, id uint, request model.UpdateProjectRequest) error
	UpdateStatus(ctx context.Context, id uint, status string) error
	GetDetail(ctx context.Context, id uint) (project model.Project, err error)
	GetAll(
		filter model.FilterProject,
		pg *utils.PaginateQueryOffset,
	) (response model.ListProjectResponse, err error)
}

type ProjectServiceImpl struct {
	Repository     repository.ProjectRepository
	TeamRepo       tm.TeamRepository
	TeamMemberRepo tm.TeamMemberRepository
	EventRepo      eve.EventRepository
}

func NewProjectService(
	repository repository.ProjectRepository,
	teamRepo tm.TeamRepository,
	teamMemberRepo tm.TeamMemberRepository,
	eventRepo eve.EventRepository,
) ProjectService {
	return &ProjectServiceImpl{
		Repository:     repository,
		TeamRepo:       teamRepo,
		TeamMemberRepo: teamMemberRepo,
		EventRepo:      eventRepo,
	}
}

func (service *ProjectServiceImpl) Create(ctx context.Context, request model.CreateProjectRequest) error {
	authenticatedUser := ctx.Value("user").(um.User)
	event, err := service.EventRepo.FindOne(request.EventID)
	if err != nil {
		return err
	}

	if event.Status != constants.EventRunning {
		return e.ErrEventNotRunning
	}

	//TODO: check timeline

	team, err := service.TeamRepo.FindOne(request.TeamID)
	if err != nil {
		return err
	}

	if team.ParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	project, err := service.Repository.FindByEventIDAndTeamID(request.EventID, request.TeamID)
	if err != nil && err != e.ErrDataNotFound {
		return err
	}

	if project.ID != 0 {
		return e.ErrTeamAlreadyHasProject
	}

	var submittedAt *time.Time
	if request.Status == constants.ProjectStatusSubmitted {
		submittedAt = helper.ReferTime(time.Now())
	}

	newProject := model.Project{
		BaseEntity:    builder.BuildBaseEntity(ctx, true, nil),
		TeamID:        request.TeamID,
		EventID:       request.EventID,
		Name:          request.Name,
		Thumbnail:     request.Thumbnail,
		ElevatorPitch: request.ElevatorPitch,
		Story:         request.Story,
		Video:         request.Video,
		Status:        request.Status,
		SubmittedAt:   submittedAt,
	}

	for k := range request.BuiltWith {
		newProject.BuiltWith = append(newProject.BuiltWith, model.ProjectTechnology{
			TechnologyID: request.BuiltWith[k],
		})
	}
	for k := range request.SiteLinks {
		newProject.SiteLinks = append(newProject.SiteLinks, model.ProjectSiteLink{
			BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
			Link:       request.SiteLinks[k],
		})
	}
	for k := range request.Images {
		newProject.Images = append(newProject.Images, model.ProjectImage{
			BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
			Image:      request.Images[k],
		})
	}

	if err = service.Repository.Create(newProject); err != nil {
		return err
	}

	return nil
}

func (service *ProjectServiceImpl) Update(ctx context.Context, id uint, request model.UpdateProjectRequest) error {
	authenticatedUser := ctx.Value("user").(um.User)
	project, err := service.Repository.FindOne(id)
	if err != nil {
		return err
	}

	event, err := service.EventRepo.FindOne(project.EventID)
	if err != nil {
		return err
	}

	if event.Status != constants.EventRunning {
		return e.ErrEventNotRunning
	}

	//TODO: check timeline

	if project.Status != constants.ProjectStatusDraft {
		return e.ErrProjectStatusShouldBeDraft
	}

	if project.Team.ParticipantID != authenticatedUser.Participant.ID {
		return e.ErrForbidden
	}

	project.Name = request.Name
	project.Thumbnail = request.Thumbnail
	project.ElevatorPitch = request.ElevatorPitch
	project.Story = request.Story
	project.Video = request.Video
	project.Status = request.Status
	project.UpdatedAt = time.Now()
	project.UpdatedBy = authenticatedUser.Email
	if request.Status == constants.ProjectStatusSubmitted {
		project.SubmittedAt = helper.ReferTime(time.Now())
	}
	project.BuiltWith = nil
	project.SiteLinks = nil
	project.Images = nil

	for k := range request.BuiltWith {
		project.BuiltWith = append(project.BuiltWith, model.ProjectTechnology{
			ProjectID:    project.ID,
			TechnologyID: request.BuiltWith[k],
		})
	}
	for k := range request.SiteLinks {
		project.SiteLinks = append(project.SiteLinks, model.ProjectSiteLink{
			BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
			ProjectID:  project.ID,
			Link:       request.SiteLinks[k],
		})
	}
	for k := range request.Images {
		project.Images = append(project.Images, model.ProjectImage{
			BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
			ProjectID:  project.ID,
			Image:      request.Images[k],
		})
	}

	var removedBuiltWith []model.ProjectTechnology
	if len(request.RemovedBuiltWith) > 0 {
		for k := range request.RemovedBuiltWith {
			removedBuiltWith = append(removedBuiltWith, model.ProjectTechnology{
				ProjectID:    project.ID,
				TechnologyID: request.RemovedBuiltWith[k],
			})
		}
	}

	if err = service.Repository.Update(id, model.UpdateProjectModel{
		Project:          project,
		RemovedBuiltWith: removedBuiltWith,
		RemovedSiteLinks: request.RemovedSiteLinks,
		RemovedImages:    request.RemovedImages,
	}); err != nil {
		return err
	}
	return nil
}

func (service *ProjectServiceImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	if status != constants.ProjectStatusInactive {
		return e.ErrInvalidStatus
	}

	authenticatedUser := ctx.Value("user").(um.User)
	project, err := service.Repository.FindOne(id)
	if err != nil {
		return err
	}

	event, err := service.EventRepo.FindOne(project.EventID)
	if err != nil {
		return err
	}

	if event.Status != constants.EventRunning {
		return e.ErrEventNotRunning
	}

	project.Status = status
	project.UpdatedAt = time.Now()
	project.UpdatedBy = authenticatedUser.Email
	if err = service.Repository.UpdateStatus(id, project); err != nil {
		return err
	}

	return nil
}

func (service *ProjectServiceImpl) GetDetail(ctx context.Context, id uint) (project model.Project, err error) {
	authenticatedUser := ctx.Value("user").(um.User)
	if project, err = service.Repository.FindOne(id); err != nil {
		return
	}

	if authenticatedUser.UserRole.Name == constants.UserParticipant {
		_, err = service.TeamMemberRepo.FindByParticipantIDAndTeamID(authenticatedUser.Participant.ID, project.TeamID)
		if err != nil && err != e.ErrDataNotFound {
			return
		} else if err != nil && err == e.ErrDataNotFound {
			err = e.ErrForbidden
			return
		}
	}

	return
}

func (service *ProjectServiceImpl) GetAll(
	filter model.FilterProject,
	pg *utils.PaginateQueryOffset,
) (response model.ListProjectResponse, err error) {
	response.Projects, response.TotalItem, response.TotalPage, err = service.Repository.FindAll(filter, pg)
	if err != nil {
		return
	}
	return
}
