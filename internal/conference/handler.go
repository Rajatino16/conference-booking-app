package conference

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, repo Repository) {
	h := NewHandler(repo)
	group := router.Group("/conference")
	{
		group.POST("", h.AddConference)
	}
}

type Handler struct {
	service Service
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		service: NewService(repo),
	}
}

func (h *Handler) AddConference(c *gin.Context) {
	var req AddConferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddConference(req); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"conference created": true})
}
