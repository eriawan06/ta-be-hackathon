package model

import "be-sagara-hackathon/src/utils/common"

type Skill struct {
	common.BaseEntity
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	IsActive bool   `gorm:"not null;default:true" json:"is_active"`
}

type SkillRequest struct {
	Name string `json:"name" validate:"required,max=20"`
}

type UpdateSkillRequest struct {
	SkillRequest
	IsActive bool `json:"is_active" validate:"omitempty"`
}

type FilterSkill struct {
	Name   string
	Status string
}

type SkillLite struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type ListSkillResponse struct {
	Skills    []SkillLite `json:"skills"`
	TotalPage int64       `json:"total_page"`
	TotalItem int64       `json:"total_item"`
}
