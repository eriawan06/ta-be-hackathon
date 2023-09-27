package model

import (
	evm "be-sagara-hackathon/src/modules/event/model"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
)

type ProjectAssessment struct {
	common.BaseEntity
	JudgeID    uint                        `gorm:"not null"`
	Judge      um.User                     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"judge"`
	ProjectID  uint                        `gorm:"not null"`
	Project    Project                     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	CriteriaID uint                        `gorm:"not null"`
	Criteria   evm.EventAssessmentCriteria `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"criteria"`
	Score      uint                        `gorm:"not null"`
}

type ProjectAssessmentRequest struct {
	CriteriaID uint `json:"criteria_id" validate:"required"`
	Score      uint `json:"score" validate:"required"`
}

type CreateBatchProjectAssessmentRequest struct {
	Assessments []ProjectAssessmentRequest `json:"assessments"`
}

type GetByProjectIDResponse struct {
	JudgeID     uint                `json:"judge_id"`
	JudgeName   string              `json:"judge_name"`
	Assessments []ProjectAssessment `json:"assessments"`
}
