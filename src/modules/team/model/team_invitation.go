package model

import (
	evm "be-sagara-hackathon/src/modules/event/model"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type TeamInvitation struct {
	common.BaseEntity
	Code            string         `gorm:"uniqueIndex;type:varchar(255);not nul"`
	EventID         uint           `gorm:"not null;"`
	Event           evm.Event      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	TeamID          uint           `gorm:"not null;"`
	Team            Team           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ToParticipantID uint           `gorm:"not null;"`
	ToParticipant   um.Participant `gorm:"foreignKey:ToParticipantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status          string         `gorm:"type:varchar(10);not null;"` //sent,accepted,rejected
	Note            *string        `gorm:"type:text"`
	ProceedAt       *time.Time     `gorm:"null"`
}

type CreateInvitationRequest struct {
	EventID         uint    `json:"event_id" validate:"required"`
	TeamID          uint    `json:"team_id" validate:"required"`
	ToParticipantID uint    `json:"to_participant_id" validate:"required"`
	Note            *string `json:"note" validate:"omitempty"`
}

type UpdateInvitationRequest struct {
	Note *string `json:"note" validate:"omitempty"`
}

type UpdateStatusInvitationRequest struct {
	Status string `json:"status" validate:"required,oneof=accepted rejected"`
}

type FilterInvitation struct {
	EventID       uint
	ParticipantID uint
	Status        string
}

type InvitationLite struct {
	ID          uint    `json:"id"`
	Code        string  `json:"code"`
	TeamID      uint    `json:"team_id"`
	TeamCode    string  `json:"team_code"`
	Name        string  `json:"name"`
	NumOfMember uint    `json:"num_of_member"`
	Avatar      *string `json:"avatar"`
	Status      string  `json:"status"`
}

type GetListInvitationResponse struct {
	Invitations []InvitationLite `json:"invitations"`
	TotalPage   int64            `json:"total_page"`
	TotalItem   int64            `json:"total_item"`
}

type InvitationDetail struct {
	ID              uint                    `json:"id"`
	Code            string                  `json:"code"`
	Note            *string                 `json:"note"`
	ToParticipantID uint                    `json:"-"`
	TeamID          uint                    `json:"team_id"`
	TeamCode        string                  `json:"team_code"`
	Name            string                  `json:"name"`
	NumOfMember     uint                    `json:"num_of_member"`
	Description     *string                 `json:"description"`
	Avatar          *string                 `json:"avatar"`
	Status          string                  `json:"status"`
	TeamMembers     []TeamMemberParticipant `json:"members" gorm:"-"`
}

type TeamInvitationList struct {
	um.ParticipantSearch
	InvitationID   uint   `json:"invitation_id"`
	InvitationCode string `json:"invitation_code"`
	Status         string `json:"invitation_status"`
}

type TeamInvitationDetail struct {
	TeamInvitationList
	CreatedAt *time.Time `json:"invitation_created_at"`
	UpdatedAt *time.Time `json:"invitation_updated_at"`
	ProceedAt *time.Time `json:"invitation_proceed_at"`
	Note      *string    `json:"invitation_note"`
	TeamID    uint       `json:"-"`
}
