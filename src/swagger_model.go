package src

import (
	ed "be-sagara-hackathon/src/modules/event/model"
	upd "be-sagara-hackathon/src/modules/general/upload/model"
	spm "be-sagara-hackathon/src/modules/master-data/speciality/model"
	pym "be-sagara-hackathon/src/modules/payment/model"
)

// Token Response
type Token struct {
	Token string `json:"token"`
}

// Error Type
type Error struct {
	Error []interface{}
}

// BaseSuccess Base Success Response
type BaseSuccess struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Success"`
	Success bool   `json:"success" example:"true"`
}

// BaseFailure Base Failure Response
type BaseFailure struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Failure"`
	Success bool   `json:"success" example:"false"`
	Error   Error  `json:"error"`
}

type AuthSuccess struct {
	BaseSuccess
	Data Token `json:"data"`
}

type ListSpecialitySuccess struct {
	BaseSuccess
	Data []spm.Speciality
}

//type GetUserProfileSuccess struct {
//	BaseSuccess
//	Data ud.UserProfileResponse
//}
//
//type GetParticipantProfileSuccess struct {
//	BaseSuccess
//	Data ud.ParticipantProfileResponse
//}
//
//type GetListMentorSuccess struct {
//	BaseSuccess
//	Data []model.MentorResponse
//}
//
//type GetMentorSuccess struct {
//	BaseSuccess
//	Data model.MentorResponse
//}
//
//type GetListJudgeSuccess struct {
//	BaseSuccess
//	Data []model.JudgeResponse
//}
//
//type GetJudgeSuccess struct {
//	BaseSuccess
//	Data model.JudgeResponse
//}

type GetListEventsSuccess struct {
	BaseSuccess
	Data []ed.EventResponse
}

type GetEventSuccess struct {
	BaseSuccess
	Data ed.EventResponse
}

type GetListEventJudgeSuccess struct {
	BaseSuccess
	Data []ed.EventJudgeLite
}

type GetEventJudgeSuccess struct {
	BaseSuccess
	Data ed.EventJudge
}

type GetListEventMentorSuccess struct {
	BaseSuccess
	Data []ed.EventMentorLite
}

type GetEventMentorSuccess struct {
	BaseSuccess
	Data ed.EventMentor
}

type GetListEventTimelineSuccess struct {
	BaseSuccess
	Data []ed.EventTimeline
}

type GetEventTimelineSuccess struct {
	BaseSuccess
	Data ed.EventTimeline
}

type GetListPaymentMethodsSuccess struct {
	BaseSuccess
	Data pym.ListPaymentMethodResponse
}

type GetPaymentMethodDetailSuccess struct {
	BaseSuccess
	Data pym.PaymentMethod
}

type GetListInvoicesFullSuccess struct {
	BaseSuccess
	Data pym.ListInvoiceResponse
}

type GetInvoiceFullSuccess struct {
	BaseSuccess
	Data pym.InvoiceFull
}

type GetListPaymentSuccess struct {
	BaseSuccess
	Data []pym.PaymentLite
}

type GetPaymentDetailSuccess struct {
	BaseSuccess
	Data pym.PaymentDetail
}

//type GetListTeamSuccess struct {
//	BaseSuccess
//	Data []td.TeamLiteResponse
//}
//
//type GetTeamDetailSuccess struct {
//	BaseSuccess
//	Data td.TeamDetailResponse
//}
//
//type GetListProjectSuccess struct {
//	BaseSuccess
//	Data []pre.ProjectLite
//}
//
//type GetProjectDetailSuccess struct {
//	BaseSuccess
//	Data []pre.ProjectDetail
//}

type UploadRequest struct {
	File      string `json:"file" binding:"required"`
	Path      string `json:"path" binding:"required"`
	Overwrite bool   `json:"overwrite"`
	PrevFile  string `json:"previous_file"`
}

type UploadFileSuccess struct {
	BaseSuccess
	upd.UploadResponse
}
