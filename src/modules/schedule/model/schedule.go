package model

import (
	evm "be-sagara-hackathon/src/modules/event/model"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type Schedule struct {
	common.BaseEntity
	EventID  uint      `gorm:"not null"`
	Event    evm.Event `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MentorID uint      `gorm:"not null"`
	Mentor   um.User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Title    string    `gorm:"type:varchar(255);not null"`
	HeldOn   time.Time `gorm:"not null"`
}

type ScheduleTeam struct {
	ScheduleID uint `gorm:"primaryKey;autoIncrement:false" json:"schedule_id" validate:"required"`
	TeamID     uint `gorm:"primaryKey;autoIncrement:false" json:"team_id" validate:"required"`
}

type ScheduleRequest struct {
	EventID  uint   `json:"event_id" validate:"required"`
	MentorID uint   `json:"mentor_id" validate:"required"`
	Title    string `json:"title" validate:"required"`
	HeldOn   string `json:"held_on" validate:"required"`
}

type FilterSchedule struct {
	EventID  uint
	HeldOn   string
	Search   string // title, mentor name
	MentorID uint
}

type ScheduleLite struct {
	ID         uint      `json:"id"`
	EventID    uint      `json:"event_id"`
	MentorName string    `json:"mentor_name"`
	Title      string    `json:"title"`
	HeldOn     time.Time `json:"held_on"`
}

type ListScheduleResponse struct {
	Schedules []ScheduleLite `json:"schedules"`
	TotalPage int64          `json:"total_page"`
	TotalItem int64          `json:"total_item"`
}

type ScheduleDetail struct {
	ScheduleLite
	MentorOccupation  string  `json:"mentor_occupation"`
	MentorInstitution string  `json:"mentor_institution"`
	MentorAvatar      *string `json:"mentor_avatar"`
}

type ScheduleLite2 struct {
	ID     uint      `json:"id"`
	Title  string    `json:"title"`
	HeldOn time.Time `json:"held_on"`
}
