package model

import (
	eve "be-sagara-hackathon/src/modules/event/model"
	ue "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type Invoice struct {
	common.BaseEntity
	EventID       uint           `gorm:"not null"`
	Event         eve.Event      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	ParticipantID uint           `gorm:"not null"`
	Participant   ue.Participant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	InvoiceNumber string         `gorm:"type:varchar(50);uniqueIndex;not null"`
	Amount        uint64         `gorm:"not null"`
	Status        string         `gorm:"type:varchar(15)"`
	PaidAmount    uint64         `gorm:"default:0"`
	ApprovedAt    *time.Time     `gorm:"null"`
	ApprovedBy    *string        `gorm:"null;type:varchar(36)"`
}

type FilterInvoice struct {
	Search  string
	EventID uint
	Status  string
}

type InvoiceLite struct {
	ID              uint   `json:"id"`
	EventID         uint   `json:"event_id"`
	InvoiceNumber   string `json:"invoice_number"`
	ParticipantName string `json:"participant_name"`
	Amount          uint64 `json:"amount"`
	Status          string `json:"status"`
}

type ListInvoiceResponse struct {
	Invoices  []InvoiceLite `json:"invoices"`
	TotalPage int64         `json:"total_page"`
	TotalItem int64         `json:"total_item"`
}

type InvoiceFull struct {
	ID               uint       `json:"id"`
	EventID          uint       `json:"event_id"`
	EventName        string     `json:"event_name"`
	InvoiceNumber    string     `json:"invoice_number"`
	ParticipantID    uint       `json:"participant_id"`
	ParticipantName  string     `json:"participant_name"`
	ParticipantEmail string     `json:"participant_email"`
	ParticipantPhone string     `json:"participant_phone"`
	Status           string     `json:"status"`
	Amount           uint64     `json:"amount"`
	PaidAmount       uint64     `json:"paid_amount"`
	ApprovedAt       *time.Time `json:"approved_at"`
	ApprovedBy       *string    `json:"approved_by"`
}
