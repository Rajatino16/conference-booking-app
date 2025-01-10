package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, repo Repository) {
	h := NewHandler(repo)
	group := router.Group("/user")
	{
		group.POST("", h.AddUser)
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

func (h *Handler) AddUser(c *gin.Context) {
	var req AddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddUser(req); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
