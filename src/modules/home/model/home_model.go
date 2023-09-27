package model

import (
	evm "be-sagara-hackathon/src/modules/event/model"
)

type HomeResponse struct {
	Event          evm.EventResponse     `json:"event_info"`
	EventJudges    []evm.EventJudgeLite  `json:"event_judges"`
	EventMentors   []evm.EventMentorLite `json:"event_mentors"`
	EventTimeline  []evm.EventTimeline   `json:"event_timelines"`
	EventFaqs      []evm.EventFaq        `json:"event_faqs"`
	EventCompanies []evm.EventCompany    `json:"event_companies"`
}
