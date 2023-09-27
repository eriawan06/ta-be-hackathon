package model

import (
	evm "be-sagara-hackathon/src/modules/event/model"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type TeamRequest struct {
	common.BaseEntity
	Code          string         `gorm:"uniqueIndex;type:varchar(255);not nul" json:"code"`
	EventID       uint           `gorm:"not null;" json:"event_id"`
	Event         evm.Event      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	TeamID        uint           `gorm:"not null;" json:"team_id"`
	Team          Team           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"team"`
	ParticipantID uint           `gorm:"not null;" json:"participant_id"`
	Participant   um.Participant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Status        string         `gorm:"type:varchar(10);not null;" json:"status"` //sent,accepted,rejected
	Note          *string        `gorm:"type:text;null" json:"note"`
	ProceedBy     *string        `gorm:"type:varchar(255);null" json:"proceed_by"`
	ProceedAt     *time.Time     `gorm:"null" json:"proceed_at"`
}

type CreateRequestJoinTeam struct {
	EventID uint    `json:"event_id" validate:"required"`
	TeamID  uint    `json:"team_id" validate:"required"`
	Note    *string `json:"note" validate:"omitempty"`
}

type UpdateRequestJoinTeam struct {
	Note *string `json:"note" validate:"omitempty"`
}

type UpdateStatusRequestJoinTeam struct {
	Status string `json:"status" validate:"required,oneof=accepted rejected"`
}

type TeamRequestList struct {
	um.ParticipantSearch
	RequestID   uint   `json:"request_id"`
	RequestCode string `json:"request_code"`
	Status      string `json:"request_status"`
}

type TeamRequestDetail struct {
	TeamRequestList
	CreatedAt *time.Time `json:"request_created_at"`
	UpdatedAt *time.Time `json:"request_updated_at"`
	ProceedAt *time.Time `json:"request_proceed_at"`
	ProceedBy *string    `json:"request_proceed_by"`
	Note      *string    `json:"request_note"`
	TeamID    uint       `json:"-"`
}
