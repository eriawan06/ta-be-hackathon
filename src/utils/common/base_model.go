package common

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time  `gorm:"not null;autoCreateTime" json:"created_at"`
	CreatedBy string     `gorm:"type:varchar(36);null;default:NULL" json:"created_by"`
	UpdatedAt time.Time  `gorm:"not null;autoUpdateTime" json:"updated_at"`
	UpdatedBy string     `gorm:"type:varchar(36);null;default:NULL" json:"updated_by"`
	DeletedAt *time.Time `gorm:"default:NULL" json:"deleted_at"`
	DeletedBy *string    `gorm:"type:varchar(36);null;default:NULL" json:"deleted_by"`
}

type BaseDtoResponse struct {
	Id        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy sql.NullString `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy sql.NullString `json:"updated_by"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	DeletedBy sql.NullString `json:"deleted_by"`
}
