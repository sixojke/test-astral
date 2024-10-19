package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sixojke/test-astral/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Init(router *gin.RouterGroup) {
	register := router.Group("/register")
	{
		register.GET("", h.registerUser)
	}
}
