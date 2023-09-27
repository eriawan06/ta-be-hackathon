package model

import (
	pye "be-sagara-hackathon/src/modules/payment/model"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type VerificationCode struct {
	common.BaseEntity
	UserID   uint
	User     um.User
	Code     string `gorm:"uniqueIndex;type:varchar(255)"`
	Type     string `gorm:"type:varchar(20)"`
	ExpireAt *time.Time
}

type RegisterModel struct {
	User          um.User
	LatestEventID uint
	Invoice       pye.Invoice
	Verification  *VerificationCode
}

type RegisterRequest struct {
	FullName        string `json:"full_name" validate:"required"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterByGoogleRequest struct {
	IdToken     string `json:"id_token" validate:"required"`
	FullName    string `json:"full_name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type LoginByGoogleRequest struct {
	IdToken string `json:"id_token" validate:"required"`
}

type VerifyEmailRequest struct {
	//Email string `json:"email" validate:"required"`
	Code string `json:"verification_code" validate:"required"`
}

type ValidateVerificationCodeRequest struct {
	Code string `json:"verification_code" validate:"required"`
}

type SendVerificationCodeRequest struct {
	Type  string `json:"type" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type ForgotPasswordRequest struct {
	Code        string `json:"verification_code" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" validate:"required"`
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type UserClaims struct {
	Authorized    bool    `json:"authorized"`
	UserId        uint    `json:"user_id"`
	ParticipantId *uint   `json:"participant_id"`
	StatusTeam    *string `json:"status_team"`
	Email         string  `json:"email"`
	RoleId        uint    `json:"role_id"`
	RoleName      string  `json:"role_name"`
	Expired       int64   `json:"expired"`
}

type AuthResponse struct {
	Token string          `json:"token"`
	User  um.UserResponse `json:"user"`
}
