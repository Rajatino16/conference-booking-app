package booking

import (
	"errors"
	"time"

	"conference-booking/internal/conference"
	"conference-booking/internal/user"

	"github.com/google/uuid"
)

var (
	ErrInvalidAction   = errors.New("action not allowed")
	ErrSlotUnavailable = errors.New("no slots available")
	ErrBookingConflict = errors.New("user already has a confirmed booking")
	ErrWaitlistExpired = errors.New("waitlist confirmation expired")
)

type Service interface {
	BookConference(req BookConferenceRequest) (string, error)
	ConfirmWaitlistBooking(bookingID string) error
	CancelBooking(bookingID string) error
	GetBookingStatus(bookingID string) (*BookingStatus, error)
}

type service struct {
	confRepo    conference.Repository
	userRepo    user.Repository
	bookingRepo Repository
}

func NewService(confRepo conference.Repository, userRepo user.Repository, bookingRepo Repository) Service {
	return &service{confRepo: confRepo, userRepo: userRepo, bookingRepo: bookingRepo}
}

func (s *service) BookConference(req BookConferenceRequest) (string, error) {
	conf, err := s.confRepo.FindByName(req.ConferenceName)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return "", err
	}

	if conf.AvailableSlots > 0 {
		// Check for overlapping bookings
		bookingID := uuid.New().String()
		booking := &Booking{
			ID:           bookingID,
			UserID:       user.ID,
			ConferenceID: conf.Name,
			Status:       "Confirmed",
		}
		s.bookingRepo.Create(booking)

		conf.AvailableSlots--
		return bookingID, nil
	}

	// Add to waitlist
	bookingID := uuid.New().String()
	waitlistUntil := time.Now().Add(1 * time.Hour)
	booking := &Booking{
		ID:            bookingID,
		UserID:        user.ID,
		ConferenceID:  conf.Name,
		Status:        "Waitlisted",
		WaitlistUntil: &waitlistUntil,
	}
	s.bookingRepo.Create(booking)
	return bookingID, nil
}

func (s *service) ConfirmWaitlistBooking(bookingID string) error {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return err
	}

	if booking.Status != "Waitlisted" || booking.WaitlistUntil.Before(time.Now()) {
		return ErrWaitlistExpired
	}

	booking.Status = "Confirmed"
	s.bookingRepo.Update(booking)
	return nil
}

func (s *service) CancelBooking(bookingID string) error {
	return s.bookingRepo.Cancel(bookingID)
}

func (s *service) GetBookingStatus(bookingID string) (*BookingStatus, error) {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return nil, err
	}

	return &BookingStatus{
		Status:        booking.Status,
		WaitlistUntil: booking.WaitlistUntil,
	}, nil
}
