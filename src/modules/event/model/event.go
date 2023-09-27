package model

import (
	ue "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type Event struct {
	common.BaseEntity
	UserID         uint            `gorm:"not null" json:"user_id"`
	User           *ue.User        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Name           string          `gorm:"type:varchar(255);not null" json:"name"`
	StartDate      time.Time       `gorm:"not null" json:"start_date"`
	EndDate        time.Time       `gorm:"not null" json:"end_date"`
	RegFee         uint64          `gorm:"not null;default:0" json:"reg_fee"`
	PaymentDueDate time.Time       `gorm:"not null" json:"payment_due_date"`
	TeamMinMember  uint            `gorm:"not null;default:0" json:"team_min_member"`
	TeamMaxMember  uint            `gorm:"not null;default:0" json:"team_max_member"`
	Description    string          `gorm:"type:text" json:"description"`
	Status         string          `gorm:"type:varchar(15);default:created" json:"status"` //status : created [PASS], approved [PASS], rejected [PASS], running, finished, inactive
	Mentors        []EventMentor   `json:"mentors" json:"mentors"`
	Judges         []EventJudge    `json:"judges" json:"judges"`
	Timelines      []EventTimeline `json:"timelines" json:"timelines"`
	Companies      []EventCompany  `json:"companies" json:"companies"`
	Rules          []EventRule     `json:"rules" json:"rules"`
	FAQs           []EventFaq      `json:"faqs" json:"faqs"`
}

type CreateEventRequest struct {
	Name           string `json:"name"  validate:"required"`
	StartDate      string `json:"start_date" validate:"required"`
	EndDate        string `json:"end_date" validate:"required"`
	RegFee         uint64 `json:"reg_fee"  validate:"required"`
	PaymentDueDate string `json:"payment_due_date"  validate:"required"`
	TeamMinMember  uint   `json:"team_min_member" validate:"required"`
	TeamMaxMember  uint   `json:"team_max_member" validate:"required"`
	Description    string `json:"description"  validate:"required"`
}

type UpdateEventRequest struct {
	CreateEventRequest
	Status string `json:"status" validate:"required"`
}

type FilterEvent struct {
	StartDate string
	EndDate   string
	Status    string
	Search    string
}

type EventLite struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	RegFee    uint64    `json:"reg_fee"`
	Status    string    `json:"status"`
}

type ListEventResponse struct {
	Events    []EventLite `json:"events"`
	TotalPage int64       `json:"total_page"`
	TotalItem int64       `json:"total_item"`
}

type EventResponse struct {
	Id             uint      `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `json:"status"`
	RegFee         uint64    `json:"reg_fee"`
	PaymentDueDate time.Time `json:"payment_due_date" `
	TeamMinMember  uint      `json:"team_min_member"`
	TeamMaxMember  uint      `json:"team_max_member"`
}
