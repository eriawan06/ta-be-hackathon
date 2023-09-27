package model

import (
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
)

type EventJudge struct {
	common.BaseEntity
	EventID uint    `gorm:"not null" json:"event_id"`
	Event   Event   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	JudgeID uint    `gorm:"not null" json:"judge_id"`
	Judge   um.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"judge"`
}

type EventJudgeRequest struct {
	EventID uint `json:"event_id" validate:"required"`
	JudgeID uint `json:"judge_id" validate:"required"`
}

type FilterEventJudge struct {
	EventID uint
}

type EventJudgeLite struct {
	ID               uint   `json:"id"`
	EventID          uint   `json:"event_id"`
	JudgeID          uint   `json:"judge_id"`
	JudgeName        string `json:"name"`
	JudgeOccupation  string `json:"occupation"`
	JudgeInstitution string `json:"institution"`
	JudgeAvatar      string `json:"avatar"`
}
