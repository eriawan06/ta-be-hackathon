package model

import (
	"be-sagara-hackathon/src/utils/common"
)

type EventAssessmentCriteria struct {
	common.BaseEntity
	EventID       uint   `gorm:"not null" json:"event_id"`
	Event         Event  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Criteria      string `gorm:"type:varchar(255);not null" json:"criteria"`
	PercentageVal uint   `gorm:"not null" json:"percentage_val"`
	ScoreStart    uint   `gorm:"not null" json:"score_start"`
	ScoreEnd      uint   `gorm:"not null" json:"score_end"`
	IsActive      bool   `gorm:"not null;default:true"`
}

type EventAssessmentCriteriaRequest struct {
	Action        string `json:"-"`
	EventID       uint   `json:"event_id" validate:"required_if=Action create"`
	Criteria      string `json:"criteria" validate:"required"`
	PercentageVal uint   `json:"percentage_val" validate:"required"`
	ScoreStart    uint   `json:"score_start" validate:"required"`
	ScoreEnd      uint   `json:"score_end" validate:"required"`
}

type UpdateEventAssessmentCriteriaRequest struct {
	EventAssessmentCriteriaRequest
	IsActive bool `json:"is_active" validate:"omitempty"`
}

type FilterEventAssessmentCriteria struct {
	EventID uint
	Status  string
}

type ListEventAssessmentCriteriaResponse struct {
	EventCriteria []EventAssessmentCriteria `json:"criteria"`
	TotalPage     int64                     `json:"total_page"`
	TotalItem     int64                     `json:"total_item"`
}
