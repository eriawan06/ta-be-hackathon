package model

import (
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type TeamMember struct {
	common.BaseEntity
	TeamID        uint           `gorm:"not null"`
	Team          Team           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParticipantID uint           `gorm:"not null;"`
	Participant   um.Participant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	JoinedAt      time.Time      `gorm:"not null;"`
}

type CreateTeamMemberRequest struct {
	ParticipantID uint `json:"participant_id" validate:"required"`
}

type TeamMemberParticipant struct {
	ID                uint    `json:"id"`
	ParticipantID     uint    `json:"participant_id"`
	ParticipantName   string  `json:"participant_name"`
	ParticipantAvatar *string `json:"participant_avatar"`
}

type TeamMemberList struct {
	um.ParticipantSearch
	TeamMemberID uint      `json:"team_member_id"`
	IsAdmin      bool      `json:"is_admin"`
	JoinedAt     time.Time `json:"joined_at"`
}

//type TeamMemberDetailList []TeamMemberDetail
//
//func (t *TeamMemberDetailList) Value() (value driver.Value, err error) {
//	return json.Marshal(t)
//}
//
//func (t *TeamMemberDetailList) Scan(value interface{}) error {
//	b, ok := value.([]byte)
//	if !ok {
//		return errors.New("type assertion to []byte failed")
//	}
//
//	unmarshal := json.Unmarshal(b, &t)
//	return unmarshal
//}
