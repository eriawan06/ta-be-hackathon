package model

type CreateJudgeRequest struct {
	CreateUserRequest
	OccupationID uint   `json:"occupation" validate:"required"`
	Institution  string `json:"institution" validate:"required"`
	Avatar       string `json:"avatar" validate:"required"`
}

type UpdateJudgeRequest struct {
	CreateMentorRequest
	IsActive bool `json:"is_active"`
}

type JudgeLite struct {
	UserLite
}

type ListJudgeResponse struct {
	Judges    []JudgeLite `json:"judges"`
	TotalPage int64       `json:"total_page"`
	TotalItem int64       `json:"total_item"`
}
