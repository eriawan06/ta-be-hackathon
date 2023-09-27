package model

import (
	ocm "be-sagara-hackathon/src/modules/master-data/occupation/model"
	"be-sagara-hackathon/src/utils/common"
)

type User struct {
	common.BaseEntity
	UserRoleID   uint            `gorm:"not null" json:"user_role_id"`
	UserRole     *UserRole       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"user_role"`
	Name         string          `gorm:"type:varchar(255);not null" json:"name"`
	Email        string          `gorm:"type:varchar(255);not null" json:"email"`
	PhoneNumber  *string         `gorm:"type:varchar(13)" json:"phone_number"`
	Username     *string         `gorm:"type:varchar(20)" json:"username"`
	Password     *string         `gorm:"type:varchar(255)" json:"-"`
	Avatar       *string         `gorm:"type:text" json:"avatar"`
	AuthType     string          `gorm:"type:varchar(10);not null" json:"-"`
	IsActive     bool            `gorm:"default:false" json:"is_active"`
	OccupationID *uint           `gorm:"null" json:"occupation_id"`
	Occupation   *ocm.Occupation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"occupation"`
	Institution  *string         `gorm:"type:varchar(255)" json:"institution"`
	Participant  *Participant    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

type UserRole struct {
	common.BaseEntity
	Name string `gorm:"type:varchar(255);not null"`
}

type UserProfile struct {
	FullName    string  `db:"name" json:"full_name"`
	Email       string  `db:"email" json:"email"`
	PhoneNumber *string `db:"phone_number" json:"phone_number"`
	Avatar      *string `db:"avatar" json:"avatar"`
	Institution *string `db:"institution" json:"institution"`
	RoleID      uint    `db:"role_id" json:"role_id"`
	RoleName    string  `db:"role_name" json:"role_name"`
	IsActive    bool    `db:"is_active" json:"is_active"`
}

type UserResponse struct {
	Id                      uint   `json:"id"`
	FullName                string `json:"full_name"`
	Email                   string `json:"email"`
	RoleId                  uint   `json:"role_id"`
	RoleName                string `json:"role_name"`
	IsRegistrationCompleted bool   `json:"is_registration_completed"`
	PaymentStatus           string `json:"payment_status"`
}

type CreateUserRequest struct {
	Action      string `json:"-" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,max=13"`
	Password    string `json:"password" validate:"required_if=Action create,omitempty,min=6"`
	RoleID      uint   `json:"role_id" validate:"omitempty"`
}

type UpdateUserRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,max=13"`
	IsActive    bool   `json:"is_active" validate:"omitempty"`
}

type FilterUser struct {
	RoleID       int
	Status       string
	Search       string
	ExceptUserID uint
}

type UserLite struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	IsActive    bool    `json:"is_active"`
}

type ListUserResponse struct {
	Users     []UserLite `json:"users"`
	TotalPage int64      `json:"total_page"`
	TotalItem int64      `json:"total_item"`
}

type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password" validate:"required"`
	NewPassword        string `json:"new_password" validate:"required"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required"`
}
