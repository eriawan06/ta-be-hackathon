package service

import (
	evr "be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/schedule/model"
	"be-sagara-hackathon/src/modules/schedule/repository"
	tr "be-sagara-hackathon/src/modules/team/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	ur "be-sagara-hackathon/src/modules/user/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
	"context"
	"time"
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, req model.ScheduleRequest) (schedule model.Schedule, err error)
	CreateScheduleTeam(req model.ScheduleTeam) (err error)
	UpdateSchedule(ctx context.Context, id uint, req model.ScheduleRequest) (err error)
	DeleteSchedule(id uint) (err error)
	DeleteScheduleTeam(id, teamID uint) (err error)
	GetListSchedule(
		ctx context.Context,
		filter model.FilterSchedule,
		pg *utils.PaginateQueryOffset,
	) (response model.ListScheduleResponse, err error)
	GetDetailSchedule(id uint) (schedule model.ScheduleDetail, err error)
}

type ScheduleServiceImpl struct {
	Repository repository.ScheduleRepository
	EventRepo  evr.EventRepository
	UserRepo   ur.UserRepository
	TeamRepo   tr.TeamRepository
}

func NewScheduleService(
	repository repository.ScheduleRepository,
	eventRepo evr.EventRepository,
	userRepo ur.UserRepository,
	teamRepo tr.TeamRepository,
) ScheduleService {
	return &ScheduleServiceImpl{Repository: repository, EventRepo: eventRepo, UserRepo: userRepo, TeamRepo: teamRepo}
}

func (service *ScheduleServiceImpl) CreateSchedule(ctx context.Context, req model.ScheduleRequest) (schedule model.Schedule, err error) {
	event, err := service.EventRepo.FindOne(req.EventID)
	if err != nil {
		return
	}

	if event.Status != constants.EventRunning {
		err = e.ErrEventNotRunning
		return
	}

	mentor, err := service.UserRepo.FindByID(req.MentorID)
	if err != nil {
		return
	}

	if mentor.UserRole.Name != constants.UserMentor {
		err = e.ErrDataNotFound
		return
	}

	heldOn, err := helper.ParseDateTimeStringToTime(req.HeldOn)
	if err != nil {
		return
	}

	if heldOn.Before(event.StartDate) || heldOn.After(event.EndDate) {
		err = e.ErrScheduleDateNotValid
		return
	}

	if schedule, err = service.Repository.Save(model.Schedule{
		BaseEntity: builder.BuildBaseEntity(ctx, true, nil),
		EventID:    req.EventID,
		MentorID:   req.MentorID,
		Title:      req.Title,
		HeldOn:     heldOn,
	}); err != nil {
		return
	}
	return
}

func (service *ScheduleServiceImpl) CreateScheduleTeam(req model.ScheduleTeam) (err error) {
	if _, err = service.Repository.FindOne(req.ScheduleID); err != nil {
		return
	}

	if _, err = service.TeamRepo.FindOne(req.TeamID); err != nil {
		return
	}

	if err = service.Repository.SaveScheduleTeam(req); err != nil {
		return
	}

	return
}

func (service *ScheduleServiceImpl) UpdateSchedule(ctx context.Context, id uint, req model.ScheduleRequest) (err error) {
	schedule, err := service.Repository.FindOne(id)
	if err != nil {
		return
	}

	mentor, err := service.UserRepo.FindByID(req.MentorID)
	if err != nil {
		return
	}

	if mentor.UserRole.Name != constants.UserMentor {
		err = e.ErrDataNotFound
		return
	}

	heldOn, err := helper.ParseDateTimeStringToTime(req.HeldOn)
	if err != nil {
		return
	}

	event, err := service.EventRepo.FindOne(req.EventID)
	if err != nil {
		return
	}

	if heldOn.Before(event.StartDate) || heldOn.After(event.EndDate) {
		err = e.ErrScheduleDateNotValid
		return
	}

	schedule.MentorID = req.MentorID
	schedule.Title = req.Title
	schedule.HeldOn = heldOn
	schedule.UpdatedAt = time.Now()
	schedule.UpdatedBy = ctx.Value("user").(um.User).Email
	if err = service.Repository.Update(id, schedule); err != nil {
		return
	}
	return
}

func (service *ScheduleServiceImpl) DeleteSchedule(id uint) (err error) {
	if _, err = service.Repository.FindOne(id); err != nil {
		return
	}

	if err = service.Repository.Delete(id); err != nil {
		return
	}
	return
}

func (service *ScheduleServiceImpl) DeleteScheduleTeam(id, teamID uint) (err error) {
	if _, err = service.Repository.FindOne(id); err != nil {
		return
	}

	if err = service.Repository.DeleteScheduleTeam(id, teamID); err != nil {
		return
	}
	return
}

func (service *ScheduleServiceImpl) GetListSchedule(
	ctx context.Context,
	filter model.FilterSchedule,
	pg *utils.PaginateQueryOffset,
) (response model.ListScheduleResponse, err error) {
	authenticatedUser := ctx.Value("user").(um.User)
	if authenticatedUser.UserRole.Name == constants.UserMentor {
		filter.MentorID = authenticatedUser.ID
	}

	response.Schedules, response.TotalItem, response.TotalPage, err = service.Repository.Find(filter, pg)
	if err != nil {
		return
	}
	return
}

func (service *ScheduleServiceImpl) GetDetailSchedule(id uint) (schedule model.ScheduleDetail, err error) {
	if schedule, err = service.Repository.FindDetail(id); err != nil {
		return
	}
	return
}
