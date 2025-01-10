package booking

import (
	"conference-booking/pkg/errors"
	"sync"
)

type Repository interface {
	Create(booking *Booking) error
	FindByID(id string) (*Booking, error)
	FindByUserAndConference(userID, conferenceID string) (*Booking, error)
	Update(booking *Booking) error
	Cancel(bookingID string) error
	FindWaitlistForConference(conferenceID string) []*Booking
}

type inMemoryRepository struct {
	bookings map[string]*Booking
	mutex    sync.Mutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		bookings: make(map[string]*Booking),
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
