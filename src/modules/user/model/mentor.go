package model

type CreateMentorRequest struct {
	CreateUserRequest
	OccupationID uint   `json:"occupation" validate:"required"`
	Institution  string `json:"institution" validate:"required"`
	Avatar       string `json:"avatar" validate:"required"`
}

type UpdateMentorRequest struct {
	CreateMentorRequest
	IsActive bool `json:"is_active"`
}

type MentorLite struct {
	UserLite
}

type ListMentorResponse struct {
	Mentors   []MentorLite `json:"mentors"`
	TotalPage int64        `json:"total_page"`
	TotalItem int64        `json:"total_item"`
}
