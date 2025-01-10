package booking

import (
	"net/http"

	"conference-booking/internal/conference"
	"conference-booking/internal/user"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, confRepo conference.Repository, userRepo user.Repository, bookingRepo Repository) {
	h := NewHandler(confRepo, userRepo, bookingRepo)
	group := router.Group("/booking")
	{
		group.POST("", h.BookConference)
		group.POST("/waitlist/confirm", h.ConfirmWaitlistBooking)
		group.DELETE("/:id", h.CancelBooking)
		group.GET("/:id", h.GetBookingStatus)
	}
}

type Handler struct {
	service Service
}

func NewHandler(confRepo conference.Repository, userRepo user.Repository, bookingRepo Repository) *Handler {
	return &Handler{
		service: NewService(confRepo, userRepo, bookingRepo),
	}
}

func (h *Handler) BookConference(c *gin.Context) {
	var req BookConferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookingID, err := h.service.BookConference(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"booking_id": bookingID})
}

func (h *Handler) ConfirmWaitlistBooking(c *gin.Context) {
	var req ConfirmWaitlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ConfirmWaitlistBooking(req.BookingID); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CancelBooking(c *gin.Context) {
	bookingID := c.Param("id")

	if err := h.service.CancelBooking(bookingID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetBookingStatus(c *gin.Context) {
	bookingID := c.Param("id")

	status, err := h.service.GetBookingStatus(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
