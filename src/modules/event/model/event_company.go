package model

import (
	"be-sagara-hackathon/src/utils/common"
)

type EventCompany struct {
	common.BaseEntity
	EventID           uint    `gorm:"not null" json:"event_id"`
	Event             Event   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name              string  `gorm:"type:varchar(255);not null" json:"name"`
	Email             string  `gorm:"type:varchar(255);not null" json:"email"`
	PhoneNumber       string  `gorm:"type:varchar(13);not null" json:"phone_number"`
	PartnershipType   string  `gorm:"type:varchar(7);not null" json:"partnership_type"`
	SponsorshipLevel  *string `gorm:"type:varchar(20);null" json:"sponsorship_level"`
	SponsorshipAmount uint64  `gorm:"not null" json:"sponsorship_amount"`
	Logo              string  `gorm:"type:text;not null" json:"logo"`
}

type EventCompanyRequest struct {
	Action            string  `json:"-"`
	EventID           uint    `json:"event_id" validate:"required_if=Action create"`
	Name              string  `json:"name"  validate:"required"`
	Email             string  `json:"email"  validate:"required"`
	PhoneNumber       string  `json:"phone_number"  validate:"required,max=13"`
	PartnershipType   string  `json:"partnership_type" validate:"required,oneof=sponsor media"`
	SponsorshipLevel  *string `json:"sponsorship_level" validate:"omitempty,oneof=bronze silver gold platinum"`
	SponsorshipAmount uint64  `json:"sponsorship_amount" validate:"omitempty"`
	Logo              string  `json:"logo" validate:"required"`
}

type FilterEventCompany struct {
	EventID uint
}
