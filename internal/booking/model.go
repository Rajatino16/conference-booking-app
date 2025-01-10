package booking

import (
	"time"
)

type Booking struct {
	ID            string
	UserID        string
	ConferenceID  string
	Status        string
	WaitlistUntil *time.Time
}

type BookConferenceRequest struct {
	ConferenceName string `json:"conference_name"`
	UserID         string `json:"user_id"`
}

type ConfirmWaitlistRequest struct {
	BookingID string `json:"booking_id"`
}

type BookingStatus struct {
	Status        string     `json:"status"`
	WaitlistUntil *time.Time `json:"waitlist_until,omitempty"`
}
