package model

import (
	"be-sagara-hackathon/src/utils/common"
)

type PaymentMethod struct {
	common.BaseEntity
	Name          string `gorm:"type:varchar(15);not null" json:"name"`
	BankCode      string `gorm:"type:varchar(3);not null" json:"bank_code"`
	AccountNumber string `gorm:"type:varchar(255);not null" json:"account_number"`
	AccountName   string `gorm:"type:varchar(255);not null" json:"account_name"`
	IsActive      bool   `gorm:"not null" json:"is_active"`
}

type PaymentMethodRequest struct {
	Name          string `json:"name" validate:"required"`
	BankCode      string `json:"bank_code" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
}

type UpdatePaymentMethodRequest struct {
	PaymentMethodRequest
	IsActive bool `json:"is_active"`
}

type FilterPaymentMethod struct {
	Status string
}

type ListPaymentMethodResponse struct {
	Methods   []PaymentMethod `json:"methods"`
	TotalPage int64           `json:"total_page"`
	TotalItem int64           `json:"total_item"`
}
