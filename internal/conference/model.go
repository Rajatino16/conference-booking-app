package conference

import "time"

type Conference struct {
	Name           string
	StartTime      time.Time
	EndTime        time.Time
	AvailableSlots int
}

type AddConferenceRequest struct {
	Name           string    `json:"name"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	AvailableSlots int       `json:"available_slots"`
}
