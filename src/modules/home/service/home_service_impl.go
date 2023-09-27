package service

import (
	evm "be-sagara-hackathon/src/modules/event/model"
	er "be-sagara-hackathon/src/modules/event/repository"
	"be-sagara-hackathon/src/modules/home/model"
)

type HomeService interface {
	GetData() (model.HomeResponse, error)
}

type HomeServiceImpl struct {
	EventRepo         er.EventRepository
	EventJudgeRepo    er.EventJudgeRepository
	EventMentorRepo   er.EventMentorRepository
	EventTimelineRepo er.EventTimelineRepository
	EventFaqRepo      er.EventFaqRepository
	EventCompanyRepo  er.EventCompanyRepository
}

func NewHomeService(
	eventRepository er.EventRepository,
	eventJudgeRepository er.EventJudgeRepository,
	eventMentorRepository er.EventMentorRepository,
	eventTimelineRepository er.EventTimelineRepository,
	eventFaqRepository er.EventFaqRepository,
	eventCompanyRepository er.EventCompanyRepository,
) HomeService {
	return &HomeServiceImpl{
		EventRepo:         eventRepository,
		EventJudgeRepo:    eventJudgeRepository,
		EventMentorRepo:   eventMentorRepository,
		EventTimelineRepo: eventTimelineRepository,
		EventFaqRepo:      eventFaqRepository,
		EventCompanyRepo:  eventCompanyRepository,
	}
}

func (service HomeServiceImpl) GetData() (model.HomeResponse, error) {
	//get latest event
	event, err := service.EventRepo.FindLatest()
	if err != nil {
		return model.HomeResponse{}, err
	}
	eventResponse := evm.EventResponse{
		Id:             event.ID,
		Name:           event.Name,
		Description:    event.Description,
		StartDate:      event.StartDate,
		EndDate:        event.EndDate,
		Status:         event.Status,
		RegFee:         event.RegFee,
		PaymentDueDate: event.PaymentDueDate,
	}

	//get event judges
	eventJudges, err := service.EventJudgeRepo.FindManyByEventID(event.ID)
	if err != nil {
		return model.HomeResponse{}, err
	}

	//get event mentors
	eventMentors, err := service.EventMentorRepo.FindManyByEventID(event.ID)
	if err != nil {
		return model.HomeResponse{}, err
	}

	//get event timelines
	eventTimelines, err := service.EventTimelineRepo.FindManyByEventID(event.ID)
	if err != nil {
		return model.HomeResponse{}, err
	}

	//get event faqs
	eventFaqs, err := service.EventFaqRepo.FindManyByEventID(event.ID)
	if err != nil {
		return model.HomeResponse{}, err
	}

	//get event companies
	eventCompanies, err := service.EventCompanyRepo.FindManyByEventID(event.ID)
	if err != nil {
		return model.HomeResponse{}, err
	}

	//create home response
	homeResponse := model.HomeResponse{
		Event:          eventResponse,
		EventJudges:    eventJudges,
		EventMentors:   eventMentors,
		EventTimeline:  eventTimelines,
		EventFaqs:      eventFaqs,
		EventCompanies: eventCompanies,
	}

	return homeResponse, nil
}
