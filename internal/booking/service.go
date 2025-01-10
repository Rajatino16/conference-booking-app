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
	StartBookingCleanup(interval time.Duration)
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
	// Find the conference
	conf, err := s.confRepo.FindByName(req.ConferenceName)
	if err != nil {
		return "", err
	}

	// Find the user
	_, err = s.userRepo.FindByID(req.UserID)
	if err != nil {
		return "", err
	}

	// Check if the user already has an active booking for this conference
	existingBooking, err := s.bookingRepo.FindActiveBooking(req.UserID, conf.Name)
	if err == nil {
		return "", errors.New("user already has an active booking with ID: " + existingBooking.ID)
	}

	// Create a booking
	bookingID := uuid.New().String()
	if conf.AvailableSlots > 0 {
		// Create a confirmed booking
		booking := &Booking{
			ID:           bookingID,
			UserID:       req.UserID,
			ConferenceID: conf.Name,
			Status:       "Confirmed",
		}
		if err := s.bookingRepo.Create(booking); err != nil {
			return "", err
		}

		// Reduce available slots
		conf.AvailableSlots--
		if err := s.confRepo.Update(conf); err != nil {
			return "", err
		}
		return bookingID, nil
	}

	// Add to waitlist
	waitlistUntil := time.Now().Add(1 * time.Hour)
	booking := &Booking{
		ID:            bookingID,
		UserID:        req.UserID,
		ConferenceID:  conf.Name,
		Status:        "Waitlisted",
		WaitlistUntil: &waitlistUntil,
	}
	if err := s.bookingRepo.Create(booking); err != nil {
		return "", err
	}
	return bookingID, nil
}

func (s *service) ConfirmWaitlistBooking(bookingID string) error {
	// Find the booking
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return err
	}

	// Validate waitlist status and expiration
	if booking.Status != "Waitlisted" || booking.WaitlistUntil.Before(time.Now()) {
		return ErrWaitlistExpired
	}

	// Find the conference
	conf, err := s.confRepo.FindByName(booking.ConferenceID)
	if err != nil {
		return err
	}

	// Check for available slots
	if conf.AvailableSlots <= 0 {
		return ErrSlotUnavailable
	}

	// Confirm the booking
	booking.Status = "Confirmed"
	if err := s.bookingRepo.Update(booking); err != nil {
		return err
	}

	// Reduce available slots
	conf.AvailableSlots--
	if err := s.confRepo.Update(conf); err != nil {
		return err
	}

	// Remove user from overlapping waitlists
	return s.bookingRepo.RemoveOverlappingWaitlists(booking.UserID, conf.StartTime, conf.EndTime)
}

func (s *service) CancelBooking(bookingID string) error {
	// Find the booking
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return err
	}

	// Ensure the booking isn't already canceled
	if booking.Status == "Canceled" {
		return errors.New("booking already canceled")
	}

	// Find the conference
	conf, err := s.confRepo.FindByName(booking.ConferenceID)
	if err != nil {
		return err
	}

	// Cancel the booking
	booking.Status = "Canceled"
	if err := s.bookingRepo.Update(booking); err != nil {
		return err
	}

	// Handle slot reassignment for confirmed bookings
	if booking.Status == "Confirmed" {
		// Increase available slots
		conf.AvailableSlots++

		// Assign slot to the first waitlisted user
		waitlist := s.bookingRepo.FindWaitlistForConference(conf.Name)
		if len(waitlist) > 0 {
			firstWaitlisted := waitlist[0]
			firstWaitlisted.Status = "PendingConfirmation"
			until := time.Now().Add(1 * time.Hour)
			firstWaitlisted.WaitlistUntil = &until
			if err := s.bookingRepo.Update(firstWaitlisted); err != nil {
				return err
			}
		}

		// Update the conference
		if err := s.confRepo.Update(conf); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) GetBookingStatus(bookingID string) (*BookingStatus, error) {
	// Find the booking
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return nil, err
	}

	if booking.Status == "Canceled" {
		return &BookingStatus{
			Status: booking.Status,
		}, nil
	}

	if booking.Status == "Waitlisted" && booking.WaitlistUntil.Before(time.Now()) {
		return &BookingStatus{
			Status: "Expired",
		}, nil
	}

	// Return booking status
	return &BookingStatus{
		Status:        booking.Status,
		WaitlistUntil: booking.WaitlistUntil,
	}, nil
}

func (s *service) StartBookingCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval) // Wait for the specified interval
			s.cleanupBookings()
		}
	}()
}

func (s *service) cleanupBookings() {
	bookings := s.bookingRepo.GetAllBookings()

	for _, booking := range bookings {
		// Remove expired waitlisted bookings
		if booking.Status == "Waitlisted" && booking.WaitlistUntil != nil && booking.WaitlistUntil.Before(time.Now().UTC()) {
			booking.Status = "Canceled"
			s.bookingRepo.Update(booking)
			continue
		}

		// Remove confirmed bookings from overlapping waitlists
		if booking.Status == "Confirmed" {
			conf, err := s.confRepo.FindByName(booking.ConferenceID)
			if err != nil {
				continue // Skip if conference not found
			}
			_ = s.bookingRepo.RemoveOverlappingWaitlists(booking.UserID, conf.StartTime, conf.EndTime)
		}

		// Handle expired bookings based on conference timing
		conf, err := s.confRepo.FindByName(booking.ConferenceID)
		if err != nil {
			continue // Skip if conference not found
		}
		if conf.EndTime.Before(time.Now().UTC()) {
			booking.Status = "Canceled"
			s.bookingRepo.Update(booking)
		}
	}
}
