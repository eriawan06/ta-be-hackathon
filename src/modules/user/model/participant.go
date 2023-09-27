package model

import (
	regm "be-sagara-hackathon/src/modules/master-data/region/model"
	skm "be-sagara-hackathon/src/modules/master-data/skill/model"
	spem "be-sagara-hackathon/src/modules/master-data/speciality/model"
	"be-sagara-hackathon/src/utils/common"
	"time"
)

type Participant struct {
	common.BaseEntity
	UserID         uint               `gorm:"not null" json:"user_id"`
	User           *User              `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"user"`
	Gender         *string            `gorm:"type:varchar(6);null" json:"gender"`
	Birthdate      *time.Time         `gorm:"type:date;null" json:"birthdate"`
	Bio            *string            `gorm:"type:text;null" json:"bio"`
	Address        *string            `gorm:"type:text;null" json:"address"`
	ProvinceID     *uint              `gorm:"null" json:"province_id"`
	Province       *regm.RegProvince  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"province"`
	CityID         *uint              `gorm:"null" json:"city_id"`
	City           *regm.RegCity      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"city"`
	DistrictID     *uint              `gorm:"null" json:"district_id"`
	District       *regm.RegDistrict  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"district"`
	VillageID      *uint              `gorm:"null" json:"village_id"`
	Village        *regm.RegVillage   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"village"`
	LevelOfStudy   *string            `gorm:"type:varchar(20);null" json:"level_of_study"`
	School         *string            `gorm:"type:varchar(255);null" json:"school"`
	GraduationYear uint16             `gorm:"not null" json:"graduation_year"`
	Major          *string            `gorm:"type:varchar(255);null" json:"major"`
	NumOfHackathon uint               `gorm:"not null" json:"num_of_hackathon"`
	LinkPortfolio  *string            `gorm:"type:text;null" json:"link_portfolio"`
	LinkRepository *string            `gorm:"type:text;null" json:"link_repository"`
	LinkLinkedin   *string            `gorm:"type:text;null" json:"link_linkedin"`
	Resume         *string            `gorm:"type:varchar(255);null" json:"resume"`
	PaymentStatus  string             `gorm:"type:varchar(15);default:unpaid" json:"payment_status"`
	IsRegistered   bool               `gorm:"not null;default:false" json:"is_registered"`
	SpecialityID   *uint              `gorm:"null" json:"speciality_id"`
	Speciality     *spem.Speciality   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"speciality"`
	Skills         []ParticipantSkill `json:"skills"`
}

type ParticipantSkill struct {
	ParticipantID uint       `gorm:"primaryKey;autoIncrement:false" json:"-"`
	SkillID       uint       `gorm:"primaryKey;autoIncrement:false" json:"skill_id"`
	Skill         *skm.Skill `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"skill"`
}

type ParticipantSkillLite struct {
	SkillID   uint   `json:"skill_id"`
	SkillName string `json:"skill_name"`
}

type UpdateParticipantFull struct {
	UpdateParticipantProfileRequest
	UpdateParticipantEducationRequest
	UpdateParticipantPreferenceRequest
}

type UpdateParticipantProfileRequest struct {
	Action      string  `json:"action" validate:"required"`
	Name        string  `json:"name" validate:"required_if=Action update"`
	PhoneNumber string  `json:"phone_number" validate:"omitempty"`
	Avatar      *string `json:"avatar" validate:"omitempty"`
	Bio         *string `json:"bio" validate:"omitempty"`
	Birthdate   string  `json:"birthdate" validate:"required"`
	Gender      string  `json:"gender" validate:"required"`
	Address     string  `json:"address" validate:"required"`
	ProvinceID  uint    `json:"province_id" validate:"required"`
	CityID      uint    `json:"city_id" validate:"required"`
	DistrictID  uint    `json:"district_id" validate:"required"`
	VillageID   uint    `json:"village_id" validate:"required"`
}

type UpdateParticipantEducationRequest struct {
	LevelOfStudy   string  `json:"level_of_study" validate:"required,oneof=highschool undergraduate postgraduate"`
	School         string  `json:"school" validate:"required"`
	GraduationYear uint16  `json:"graduation_year" validate:"required"`
	Major          *string `json:"major"  validate:"omitempty"`
}

type UpdateParticipantPreferenceRequest struct {
	OccupationID   uint    `json:"occupation_id" validate:"required"`
	CompanyName    *string `json:"company_name" validate:"omitempty"`
	NumOfHackathon uint    `json:"num_of_hackathon" validate:"omitempty"`
	Portfolio      *string `json:"portfolio" validate:"omitempty"`
	Repository     *string `json:"repository" validate:"omitempty"`
	Linkedin       *string `json:"linkedin" validate:"omitempty"`
	Resume         *string `json:"resume" validate:"omitempty"`
	SpecialityID   uint    `json:"speciality_id" validate:"required"`
	Skills         []uint  `json:"skills" validate:"omitempty"`
	RemovedSkills  []uint  `json:"removed_skills"  validate:"omitempty"`
}

type UpdateParticipantAccountRequest struct {
	Username    string `json:"username" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type UpdateParticipant struct {
	ID            uint
	Participant   Participant
	Skills        []ParticipantSkill
	RemovedSkills []uint
}

type ParticipantLite struct {
	UserLite
	UserID uint `json:"user_id"`
}

type ListParticipantResponse struct {
	Participants []ParticipantLite `json:"participants"`
	TotalPage    int64             `json:"total_page"`
	TotalItem    int64             `json:"total_item"`
}

type FilterParticipantSearch struct {
	Search       string `json:"search" validate:"omitempty"`
	Specialities []uint `json:"specialities" validate:"omitempty"`
	Skills       []uint `json:"skills" validate:"omitempty"`
	InTeam       bool   `json:"-" validate:"omitempty"`
}

type ParticipantSearch struct {
	ID             uint                   `json:"id"`
	UserID         uint                   `json:"user_id"`
	Avatar         *string                `json:"avatar"`
	Name           string                 `json:"name"`
	SpecialityID   uint                   `json:"speciality_id"`
	SpecialityName string                 `json:"speciality_name"`
	CityID         *uint                  `json:"city_id"`
	CityName       *string                `json:"city_name"`
	Gender         string                 `json:"gender"`
	OccupationID   uint                   `json:"occupation_id"`
	OccupationName string                 `json:"occupation_name"`
	Institution    *string                `json:"institution"`
	School         *string                `json:"school"`
	Bio            *string                `json:"bio"`
	Skills         []ParticipantSkillLite `json:"skills" gorm:"-"`
}

type ListParticipantSearchResponse struct {
	Participants []ParticipantSearch `json:"participants"`
	TotalPage    int64               `json:"total_page"`
	TotalItem    int64               `json:"total_item"`
}
