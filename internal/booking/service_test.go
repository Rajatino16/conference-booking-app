package booking

// import (
// 	"conference-booking/internal/conference"
// 	"conference-booking/internal/user"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// func setupServices() (*service, *inMemoryRepository, *conference.InMemoryRepository, *user.InMemoryRepository) {
// 	confRepo := conference.NewInMemoryRepository()
// 	userRepo := user.NewInMemoryRepository()
// 	bookingRepo := NewInMemoryRepository(confRepo)
// 	service := NewService(confRepo, userRepo, bookingRepo)
// 	return service, bookingRepo, confRepo, userRepo
// }

// func TestCleanupExpiredWaitlistedBookings(t *testing.T) {
// 	service, bookingRepo, confRepo, userRepo := setupServices()

// 	// Add a conference
// 	conf := &conference.Conference{
// 		Name:          "TechConf",
// 		StartTime:     time.Now().Add(1 * time.Hour).UTC(),
// 		EndTime:       time.Now().Add(3 * time.Hour).UTC(),
// 		AvailableSlots: 0,
// 	}
// 	confRepo.Create(conf)

// 	// Add a user
// 	user := &user.User{ID: "user1"}
// 	userRepo.Create(user)

// 	wl := time.Now().Add(-1 * time.Hour).UTC()

// 	// Add an expired waitlisted booking
// 	expiredWaitlist := &Booking{
// 		ID:            "expired1",
// 		UserID:        "user1",
// 		ConferenceID:  "TechConf",
// 		Status:        "Waitlisted",
// 		WaitlistUntil: &wl,
// 	}
// 	bookingRepo.Create(expiredWaitlist)

// 	// Run the cleanup function
// 	service.cleanupBookings()

// 	// Verify the waitlisted booking is marked as canceled
// 	updatedBooking, err := bookingRepo.FindByID("expired1")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Canceled", updatedBooking.Status)
// }
