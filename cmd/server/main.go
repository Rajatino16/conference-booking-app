package main

import (
	"log"
	"time"

	"conference-booking/internal/booking"
	"conference-booking/internal/conference"
	"conference-booking/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// In-memory storage
	conferenceStore := conference.NewInMemoryRepository()
	userStore := user.NewInMemoryRepository()
	bookingStore := booking.NewInMemoryRepository(conferenceStore)

	// Initialize services
	bookingService := booking.NewService(conferenceStore, userStore, bookingStore)

	// Start cleanup goroutine (e.g., every 15 minutes)
	bookingService.StartBookingCleanup(15 * time.Minute)

	// Register routes
	conference.RegisterRoutes(router, conferenceStore)
	user.RegisterRoutes(router, userStore)
	booking.RegisterRoutes(router, conferenceStore, userStore, bookingStore)

	log.Fatal(router.Run(":8080"))
}
