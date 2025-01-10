package booking

import (
	"conference-booking/internal/conference"
	"conference-booking/pkg/errors"
	"sync"
	"time"
)

type Repository interface {
	Create(booking *Booking) error
	FindByID(id string) (*Booking, error)
	FindByUserAndConference(userID, conferenceID string) (*Booking, error)
	Update(booking *Booking) error
	Cancel(bookingID string) error
	FindWaitlistForConference(conferenceID string) []*Booking
	FindActiveBooking(userID, conferenceID string) (*Booking, error)
	RemoveOverlappingWaitlists(userID string, start, end time.Time) error
	HasOverlappingConfirmedBookings(userID string, start, end time.Time) (bool, error)
}

type inMemoryRepository struct {
	bookings       map[string]*Booking
	mutex          sync.Mutex
	conferenceRepo conference.Repository
}

func NewInMemoryRepository(confRepo conference.Repository) Repository {
	return &inMemoryRepository{
		bookings:       make(map[string]*Booking),
		conferenceRepo: confRepo,
	}
}

func (r *inMemoryRepository) Create(booking *Booking) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.bookings[booking.ID] = booking
	return nil
}

func (r *inMemoryRepository) FindByID(id string) (*Booking, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	booking, exists := r.bookings[id]
	if !exists {
		return nil, errors.ErrNotFound
	}
	return booking, nil
}

func (r *inMemoryRepository) FindByUserAndConference(userID, conferenceID string) (*Booking, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, booking := range r.bookings {
		if booking.UserID == userID && booking.ConferenceID == conferenceID {
			return booking, nil
		}
	}
	return nil, nil
}

func (r *inMemoryRepository) Update(booking *Booking) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.bookings[booking.ID] = booking
	return nil
}

func (r *inMemoryRepository) Cancel(bookingID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if booking, exists := r.bookings[bookingID]; exists {
		booking.Status = "Cancelled"
		r.bookings[bookingID] = booking
		return nil
	}
	return errors.ErrNotFound
}

func (r *inMemoryRepository) FindWaitlistForConference(conferenceID string) []*Booking {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var waitlist []*Booking
	for _, booking := range r.bookings {
		if booking.ConferenceID == conferenceID && booking.Status == "Waitlisted" {
			waitlist = append(waitlist, booking)
		}
	}
	return waitlist
}

// New Method: FindActiveBooking
func (r *inMemoryRepository) FindActiveBooking(userID, conferenceID string) (*Booking, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, booking := range r.bookings {
		if booking.UserID == userID && booking.ConferenceID == conferenceID && booking.Status != "Cancelled" && booking.Status != "Expired" {
			return booking, nil
		}
	}
	return nil, errors.ErrNotFound
}

// New Method: RemoveOverlappingWaitlists
func (r *inMemoryRepository) RemoveOverlappingWaitlists(userID string, start, end time.Time) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, booking := range r.bookings {
		if booking.UserID == userID && booking.Status == "Waitlisted" {
			// Fetch conference details using its ID
			conf, err := r.conferenceRepo.FindByName(booking.ConferenceID)
			if err != nil {
				continue // Skip if conference not found
			}

			// Check for overlapping timeframes
			if !(end.Before(conf.StartTime) || start.After(conf.EndTime)) {
				booking.Status = "Cancelled"
				r.bookings[booking.ID] = booking
			}
		}
	}
	return nil
}

func (r *inMemoryRepository) HasOverlappingConfirmedBookings(userID string, start, end time.Time) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, booking := range r.bookings {
		if booking.UserID == userID && booking.Status == "Confirmed" {
			conf, err := r.conferenceRepo.FindByName(booking.ConferenceID)
			if err != nil {
				continue // Skip if conference not found
			}
			// Check for overlapping times
			if !(end.Before(conf.StartTime) || start.After(conf.EndTime)) {
				return true, nil
			}
		}
	}
	return false, nil
}
