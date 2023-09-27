package model

import "be-sagara-hackathon/src/utils/common"

type Speciality struct {
	common.BaseEntity
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	IsActive bool   `gorm:"not null;default:true" json:"is_active"`
}

type SpecialityRequest struct {
	Name string `json:"name" validate:"required,max=20"`
}

type UpdateSpecialityRequest struct {
	SpecialityRequest
	IsActive bool `json:"is_active" validate:"omitempty"`
}

type FilterSpeciality struct {
	Name   string
	Status string
}

type SpecialityLite struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type ListSpecialityResponse struct {
	Specialities []SpecialityLite `json:"specialities"`
	TotalPage    int64            `json:"total_page"`
	TotalItem    int64            `json:"total_item"`
}
