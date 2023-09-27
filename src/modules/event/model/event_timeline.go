package model

import (
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type EventTimeline struct {
	common.BaseEntity
	EventID   uint      `gorm:"not null" json:"event_id"`
	Event     Event     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `gorm:"not null" json:"end_date"`
	Note      string    `gorm:"type:text;not null" json:"note"`
}

type EventTimelineRequest struct {
	Action    string `json:"-"`
	EventID   uint   `json:"event_id" validate:"required_if=Action create"`
	Title     string `json:"title"  validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
	Note      string `json:"note" validate:"required"`
}

type FilterEventTimeline struct {
	EventID uint
}
