package model

import (
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type Team struct {
	common.BaseEntity
	ParticipantID uint           `gorm:"not null;" json:"participant_id"`
	Participant   um.Participant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Code          string         `gorm:"type:varchar(255);not null;" json:"code"`
	Name          string         `gorm:"type:varchar(255);not null;" json:"name"`
	Description   *string        `gorm:"type:text;null" json:"description"`
	Avatar        *string        `gorm:"type:text;null" json:"avatar"`
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
}

type TeamEvent struct {
	TeamID  uint `gorm:"primaryKey;autoIncrement:false"`
	EventID uint `gorm:"primaryKey;autoIncrement:false"`
}

type CreateTeamRequest struct {
	EventID     uint    `json:"event_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description" validate:"omitempty"`
	Avatar      *string `json:"avatar" validate:"omitempty"`
}

type UpdateTeamRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description" validate:"omitempty"`
	Avatar      *string `json:"avatar" validate:"omitempty"`
}

type UpdateTeamStatusRequest struct {
	IsActive bool `json:"is_active"  validate:"omitempty"`
}

type FilterTeam struct {
	IsForParticipant         bool
	Search                   string //name, creator/participant name
	EventID                  uint
	Status                   string
	ScheduleID               uint
	TeamRequestParticipantID uint
}

type TeamLite struct {
	ID              uint    `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	EventID         uint    `json:"event_id"`
	NumOfMember     uint    `json:"num_of_member"`
	ParticipantName string  `json:"participant_name"` //creator
	IsActive        bool    `json:"is_active"`
	Avatar          *string `json:"-"`
	IsRequested     bool    `json:"-"`
}

type GetAllTeamResponse struct {
	Teams     []TeamLite `json:"teams"`
	TotalPage int64      `json:"total_page"`
	TotalItem int64      `json:"total_item"`
}

type TeamDetail struct {
	TeamLite
	CreatedAt   time.Time `json:"created_at"`
	Description *string   `json:"description"`
	Avatar      *string   `json:"avatar"`
	ProjectID   uint      `json:"project_id"`
	ProjectLink *string   `json:"project_link"`
}

type TeamByEventID struct {
	ID          uint    `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	NumOfMember uint    `json:"num_of_member"`
	Avatar      *string `json:"avatar"`
	IsRequested bool    `json:"is_requested"`
}

type GetListTeamByEventIDResponse struct {
	Teams     []TeamByEventID `json:"teams"`
	TotalPage int64           `json:"total_page"`
	TotalItem int64           `json:"total_item"`
}

type TeamDetail2 struct {
	ID               uint                    `json:"id"`
	Code             string                  `json:"code"`
	Name             string                  `json:"name"`
	NumOfMember      uint                    `json:"num_of_member"`
	Description      *string                 `json:"description"`
	Avatar           *string                 `json:"avatar"`
	IsRequested      bool                    `json:"is_requested"`
	ParticipantID    uint                    `json:"participant_id"`
	ParticipantName  string                  `json:"-"`
	ParticipantEmail string                  `json:"-"`
	ProjectID        *uint                   `json:"project_id"`
	Members          []TeamMemberParticipant `json:"members" gorm:"-"`
}
