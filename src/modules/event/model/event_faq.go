package model

import (
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type EventFaq struct {
	common.BaseEntity
	EventID  uint   `gorm:"not null" json:"event_id"`
	Event    Event  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Title    string `gorm:"type:varchar(255);not null" json:"title"`
	Note     string `gorm:"type:text;not null" json:"note"`
	IsActive bool   `gorm:"not null;default:true"`
}

type EventFaqRequest struct {
	Action  string `json:"-"`
	EventID uint   `json:"event_id" validate:"required_if=Action create"`
	Title   string `json:"title"  validate:"required"`
	Note    string `json:"note" validate:"required"`
}

type UpdateEventFaqRequest struct {
	EventFaqRequest
	IsActive bool `json:"is_active" validate:"omitempty"`
}

type FilterEventFaq struct {
	EventID uint
	Status  string
}

type EventFaqLite struct {
	ID        uint      `json:"id"`
	Title     string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

type ListEventFaqResponse struct {
	EventFaqs []EventFaqLite `json:"faqs"`
	TotalPage int64          `json:"total_page"`
	TotalItem int64          `json:"total_item"`
}
