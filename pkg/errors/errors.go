package errors

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("resource conflict")
	ErrInvalidInput    = errors.New("invalid input")
	ErrSlotUnavailable = errors.New("no slots available")
	ErrWaitlistExpired = errors.New("waitlist confirmation expired")
	ErrBookingConflict = errors.New("user already has a confirmed booking")
	ErrInvalidAction   = errors.New("action not allowed")
)
