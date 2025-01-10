package main

import (
	"log"

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
	bookingStore := booking.NewInMemoryRepository()

	// Register routes
	conference.RegisterRoutes(router, conferenceStore)
	user.RegisterRoutes(router, userStore)
	booking.RegisterRoutes(router, conferenceStore, userStore, bookingStore)

	log.Fatal(router.Run(":8080"))
}
