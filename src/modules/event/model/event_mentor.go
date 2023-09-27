package model

import (
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
)

type EventMentor struct {
	common.BaseEntity
	EventID  uint    `gorm:"not null" json:"event_id"`
	Event    Event   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	MentorID uint    `gorm:"not null" json:"mentor_id"`
	Mentor   um.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"mentor"`
}

type EventMentorRequest struct {
	EventID  uint `json:"event_id" validate:"required"`
	MentorID uint `json:"mentor_id" validate:"required"`
}

type FilterEventMentor struct {
	EventID uint
}

type EventMentorLite struct {
	ID                uint   `json:"id"`
	EventID           uint   `json:"event_id"`
	MentorID          uint   `json:"mentor_id"`
	MentorName        string `json:"name"`
	MentorOccupation  string `json:"occupation"`
	MentorInstitution string `json:"institution"`
	MentorAvatar      string `json:"avatar"`
}
