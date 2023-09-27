package service

import (
	"be-sagara-hackathon/src/modules/event/model"
	"be-sagara-hackathon/src/modules/event/repository"
	scm "be-sagara-hackathon/src/modules/schedule/model"
	scr "be-sagara-hackathon/src/modules/schedule/repository"
	tr "be-sagara-hackathon/src/modules/team/repository"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
	"context"
	"time"
)

type EventService interface {
	CreateEvent(ctx context.Context, request model.CreateEventRequest) (event *model.Event, err error)
	UpdateEvent(ctx context.Context, request model.UpdateEventRequest, eventID uint) error
	DeleteEvent(ctx context.Context, eventID uint) error
	GetListEvent(
		filter model.FilterEvent,
		pg *utils.PaginateQueryOffset,
	) (response model.ListEventResponse, err error)
	GetDetailEvent(eventID uint) (event model.Event, err error)
	GetLatestEvent() (response model.EventResponse, err error)
	GetSchedules(ctx context.Context, eventID uint) (schedules []scm.ScheduleLite2, err error)
}

type EventServiceImpl struct {
	Repository           repository.EventRepository
	EventParticipantRepo repository.EventParticipantRepository
	TeamMemberRepo       tr.TeamMemberRepository
	ScheduleRepo         scr.ScheduleRepository
}

func NewEventRepository(
	repository repository.EventRepository,
	eventParticipantRepo repository.EventParticipantRepository,
	teamMemberRepo tr.TeamMemberRepository,
	scheduleRepo scr.ScheduleRepository,
) EventService {
	return &EventServiceImpl{
		Repository:           repository,
		EventParticipantRepo: eventParticipantRepo,
		TeamMemberRepo:       teamMemberRepo,
		ScheduleRepo:         scheduleRepo,
	}
}

func (service *EventServiceImpl) CreateEvent(ctx context.Context, request model.CreateEventRequest) (event *model.Event, err error) {
	//get latest event
	latestEvent, err := service.Repository.FindLatest()
	if err != nil && err != e.ErrDataNotFound {
		return
	}

	if latestEvent.ID != 0 && latestEvent.Status == constants.EventRunning {
		err = e.ErrLatestEventRunning
		return
	}

	// Parse Event Date
	startDate, err := helper.ParseDateStringToTime(request.StartDate)
	if err != nil {
		return
	}
	endDate, err := helper.ParseDateStringToTime(request.EndDate)
	if err != nil {
		return
	}
	paymentDueDate, err := helper.ParseDateStringToTime(request.PaymentDueDate)
	if err != nil {
		return
	}
	authUser := ctx.Value("user").(um.User)
	event = &model.Event{
		BaseEntity:     builder.BuildBaseEntity(ctx, true, nil),
		UserID:         authUser.ID,
		Name:           request.Name,
		StartDate:      startDate,
		EndDate:        endDate,
		RegFee:         request.RegFee,
		PaymentDueDate: paymentDueDate,
		TeamMinMember:  request.TeamMinMember,
		TeamMaxMember:  request.TeamMaxMember,
		Description:    request.Description,
		Status:         constants.EventRunning,
	}
	if err = service.Repository.Save(event); err != nil {
		event = nil
		return
	}
	return
}

func (service *EventServiceImpl) UpdateEvent(ctx context.Context, request model.UpdateEventRequest, eventID uint) error {
	if request.Status == constants.EventRunning {
		//get latest event
		latestEvent, err := service.Repository.FindLatest()
		if err != nil && err != e.ErrDataNotFound {
			return err
		}

		if latestEvent.Status == constants.EventRunning && latestEvent.ID != eventID {
			err = e.ErrLatestEventRunning
			return err
		}
	}

	event, err := service.Repository.FindOne(eventID)
	if err != nil {
		return err
	}

	// Parse Event Date
	startDate, err := helper.ParseDateStringToTime(request.StartDate)
	if err != nil {
		return err
	}
	endDate, err := helper.ParseDateStringToTime(request.EndDate)
	if err != nil {
		return err
	}
	paymentDueDate, err := helper.ParseDateStringToTime(request.PaymentDueDate)
	if err != nil {
		return err
	}

	event.Name = request.Name
	event.StartDate = startDate
	event.EndDate = endDate
	event.RegFee = request.RegFee
	event.PaymentDueDate = paymentDueDate
	event.TeamMinMember = request.TeamMinMember
	event.TeamMaxMember = request.TeamMaxMember
	event.Description = request.Description
	event.Status = request.Status
	event.UpdatedAt = time.Now()
	event.UpdatedBy = ctx.Value("user").(um.User).Email
	if err = service.Repository.Update(event, eventID); err != nil {
		return err
	}
	return nil
}

func (service *EventServiceImpl) DeleteEvent(ctx context.Context, eventID uint) error {
	_, err := service.Repository.FindOne(eventID)
	if err != nil {
		return err
	}

	authUser := ctx.Value("user").(um.User)
	if err = service.Repository.Delete(eventID, authUser.Email); err != nil {
		return err
	}
	return nil
}

func (service *EventServiceImpl) GetListEvent(
	filter model.FilterEvent,
	pg *utils.PaginateQueryOffset,
) (response model.ListEventResponse, err error) {
	events, totalData, totalPage, err := service.Repository.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.EventLite
	for _, v := range events {
		responseData = append(responseData, model.EventLite{
			ID:        v.ID,
			Name:      v.Name,
			StartDate: v.StartDate,
			EndDate:   v.EndDate,
			RegFee:    v.RegFee,
			Status:    v.Status,
		})
	}
	response.Events = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *EventServiceImpl) GetDetailEvent(eventID uint) (event model.Event, err error) {
	if event, err = service.Repository.FindOne(eventID); err != nil {
		return
	}
	return
}

func (service *EventServiceImpl) GetLatestEvent() (response model.EventResponse, err error) {
	event, err := service.Repository.FindLatest()
	if err != nil {
		return
	}

	response = model.EventResponse{
		Id:             event.ID,
		Name:           event.Name,
		Description:    event.Description,
		StartDate:      event.StartDate,
		EndDate:        event.EndDate,
		Status:         event.Status,
		RegFee:         event.RegFee,
		PaymentDueDate: event.PaymentDueDate,
		TeamMinMember:  event.TeamMinMember,
		TeamMaxMember:  event.TeamMaxMember,
	}
	return
}

func (service *EventServiceImpl) GetSchedules(ctx context.Context, eventID uint) (schedules []scm.ScheduleLite2, err error) {
	authenticatedUser := ctx.Value("user").(um.User)
	teamMember, err := service.TeamMemberRepo.FindByParticipantID(authenticatedUser.Participant.ID)
	if err != nil {
		return
	}

	_, err = service.EventParticipantRepo.FindOneByEventIDAndParticipantID(eventID, authenticatedUser.Participant.ID)
	if err != nil && err != e.ErrDataNotFound {
		return
	} else if err != nil && err == e.ErrDataNotFound {
		err = e.ErrForbidden
		return
	}

	if schedules, err = service.ScheduleRepo.FindByEventIDAndTeamID(eventID, teamMember.TeamID); err != nil {
		return
	}
	return
}
