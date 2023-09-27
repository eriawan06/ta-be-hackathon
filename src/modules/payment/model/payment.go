package model

import (
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type Payment struct {
	common.BaseEntity
	InvoiceID         uint           `gorm:"not null"`
	Invoice           Invoice        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PaymentType       string         `gorm:"not null"` //auto, manual
	PaymentMethodID   *uint          `gorm:"null"`     //not null when payment type is manual
	PaymentMethod     *PaymentMethod `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Status            string         `gorm:"not null;default:created"` //created,proceed
	BankName          *string        `gorm:"type:varchar(10);null"`
	BankAccountName   *string        `gorm:"type:varchar(255);null"`
	BankAccountNumber *string        `gorm:"type:varchar(255);null"`
	Evidence          *string        `gorm:"type:text;null"`
	Amount            uint64         `gorm:"not null"`
	ProceedAt         *time.Time     `gorm:"null"`
	ProceedBy         *string        `gorm:"type:varchar(255);null"`
	Note              *string        `gorm:"type:text;null"`
}

// CreatePaymentRequest Payment Manual
type CreatePaymentRequest struct {
	InvoiceID       uint   `json:"invoice_id" validate:"required"`
	PaymentMethodID uint   `json:"payment_method_id" validate:"required"`
	AccountName     string `json:"account_name" validate:"required"`
	AccountNumber   string `json:"account_number" validate:"required"`
	BankName        string `json:"bank_name" validate:"required"`
	Evidence        string `json:"evidence" validate:"required"`
}

type UpdatePaymentRequest struct {
	Amount uint64  `json:"amount" validate:"required"`
	Note   *string `json:"note" validate:"omitempty"`
}

type FilterPayment struct {
	CreatedAt   string
	Status      string
	PaymentType string
	Search      string //invoice number, payment method
}

type PaymentLite struct {
	ID            uint      `json:"id"`
	InvoiceNumber string    `json:"invoice_number"`
	PaymentType   string    `json:"payment_type"`
	PaymentMethod *string   `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
	Status        string    `json:"status"`
}

type ListPaymentResponse struct {
	Payments  []PaymentLite `json:"payments"`
	TotalPage int64         `json:"total_page"`
	TotalItem int64         `json:"total_item"`
}

type PaymentDetail struct {
	ID                uint       `json:"id"`
	EventID           uint       `json:"event_id"`
	InvoiceID         uint       `json:"invoice_id"`
	InvoiceNumber     string     `json:"invoice_number"`
	ParticipantID     uint       `json:"participant_id"`
	ParticipantName   string     `json:"participant_name"`
	PaymentType       string     `json:"payment_type"`
	PaymentMethodID   *uint      `json:"payment_method_id"`
	PaymentMethodName *string    `json:"payment_method_name"`
	CreatedAt         time.Time  `json:"created_at"`
	Status            string     `json:"status"`
	BankName          *string    `json:"bank_name"`
	BankAccountName   *string    `json:"bank_account_name"`
	BankAccountNumber *string    `json:"bank_account_number"`
	Evidence          *string    `json:"evidence"`
	Amount            uint64     `json:"amount"`
	ProceedAt         *time.Time `json:"proceed_at"`
	ProceedBy         *string    `json:"proceed_by"`
	Note              *string    `json:"note"`
}
