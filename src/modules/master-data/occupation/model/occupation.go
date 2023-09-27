package model

import "be-sagara-hackathon/src/utils/common"

type Occupation struct {
	common.BaseEntity
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	IsActive bool   `gorm:"not null;default:true" json:"is_active"`
}

type OccupationRequest struct {
	Name string `json:"name" validate:"required,max=20"`
}

type UpdateOccupationRequest struct {
	OccupationRequest
	IsActive bool `json:"is_active" validate:"omitempty"`
}

type FilterOccupation struct {
	Name   string
	Status string
}

type OccupationLite struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type ListOccupationResponse struct {
	Occupations []OccupationLite `json:"occupations"`
	TotalPage   int64            `json:"total_page"`
	TotalItem   int64            `json:"total_item"`
}
