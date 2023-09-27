package model

import (
	ue "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
)

type EventParticipant struct {
	common.BaseEntity
	EventID       uint           `gorm:"not null" json:"event_id"`
	Event         Event          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	ParticipantID uint           `gorm:"not null" json:"participant_id"`
	Participant   ue.Participant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"participant"`
}
