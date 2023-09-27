package model

import "be-sagara-hackathon/src/utils/common"

type Technology struct {
	common.BaseEntity
	Name     string `gorm:"type:varchar(255);not null"`
	IsActive bool   `gorm:"not null;default:true"`
}

type TechnologyRequest struct {
	Name string `json:"name" validate:"required,max=20"`
}

type UpdateTechnologyRequest struct {
	TechnologyRequest
	IsActive bool `json:"is_active" validate:"omitempty"`
}

type FilterTechnology struct {
	Name   string
	Status string
}

type TechnologyLite struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type ListTechnologyResponse struct {
	Technologies []TechnologyLite `json:"technologies"`
	TotalPage    int64            `json:"total_page"`
	TotalItem    int64            `json:"total_item"`
}
