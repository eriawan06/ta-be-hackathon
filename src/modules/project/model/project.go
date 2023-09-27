package model

import (
	evm "be-sagara-hackathon/src/modules/event/model"
	tecm "be-sagara-hackathon/src/modules/master-data/technology/model"
	tm "be-sagara-hackathon/src/modules/team/model"
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type Project struct {
	common.BaseEntity
	TeamID        uint                `gorm:"not null;" json:"team_id"`
	Team          tm.Team             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"team"`
	EventID       uint                `gorm:"not null;" json:"event_id"`
	Event         evm.Event           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Name          string              `gorm:"type:varchar(255);not null" json:"name"`
	Thumbnail     string              `gorm:"type:text;not null" json:"thumbnail"`
	ElevatorPitch string              `gorm:"type:text;not null" json:"elevator_pitch"`
	Story         string              `gorm:"type:text;not null" json:"story"`
	Video         string              `gorm:"type:text;not null" json:"video"`
	Status        string              `gorm:"type:varchar(10);not null" json:"status"` //draft, submitted
	SubmittedAt   *time.Time          `gorm:"null" json:"submitted_at"`
	BuiltWith     []ProjectTechnology `json:"built_with"`
	SiteLinks     []ProjectSiteLink   `json:"site_links"`
	Images        []ProjectImage      `json:"images"`
}

type ProjectSiteLink struct {
	common.BaseEntity
	ProjectID uint    `gorm:"not null" json:"project_id"`
	Project   Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Link      string  `gorm:"type:text;not null" json:"link"`
}

type ProjectImage struct {
	common.BaseEntity
	ProjectID uint    `gorm:"not null" json:"project_id"`
	Project   Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Image     string  `gorm:"type:text;not null" json:"image"`
}

type ProjectTechnology struct {
	ProjectID    uint             `gorm:"primaryKey:autoIncrement:false" json:"project_id"`
	TechnologyID uint             `gorm:"primaryKey:autoIncrement:false" json:"technology_id"`
	Technology   *tecm.Technology `json:"technology"`
}

type CreateProjectRequest struct {
	Action        string   `json:"-" validate:"required,oneof=create update"`
	TeamID        uint     `json:"team_id" validate:"required_if=Action create,omitempty"`
	EventID       uint     `json:"event_id" validate:"required_if=Action create,omitempty"`
	Name          string   `json:"name" validate:"required"`
	Thumbnail     string   `json:"thumbnail" validate:"required"`
	ElevatorPitch string   `json:"elevator_pitch" validate:"required"`
	Story         string   `json:"story" validate:"required"`
	Video         string   `json:"video" validate:"required"`
	Status        string   `json:"status" validate:"required"`
	BuiltWith     []uint   `json:"built_with" validate:"required_if=Action create,omitempty"`
	SiteLinks     []string `json:"site_links" validate:"required_if=Action create,omitempty"`
	Images        []string `json:"images" validate:"required_if=Action create,omitempty"`
}

type UpdateProjectRequest struct {
	CreateProjectRequest
	RemovedBuiltWith []uint `json:"removed_built_with" validate:"omitempty"` //technology_id
	RemovedSiteLinks []uint `json:"removed_site_links" validate:"omitempty"` //project_site_link_id
	RemovedImages    []uint `json:"removed_images" validate:"omitempty"`     //project_image_id
}

type UpdateProjectModel struct {
	Project          Project
	RemovedBuiltWith []ProjectTechnology
	RemovedSiteLinks []uint
	RemovedImages    []uint
}

type FilterProject struct {
	EventID   uint
	TeamID    uint
	CreatedAt string
	Status    string
	Search    string //project name
}

type ProjectLite struct {
	ID        uint      `json:"id"`
	EventID   uint      `json:"event_id"`
	TeamID    uint      `json:"team_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ListProjectResponse struct {
	Projects  []ProjectLite `json:"projects"`
	TotalPage int64         `json:"total_page"`
	TotalItem int64         `json:"total_item"`
}
