package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sixojke/test-astral/internal/config"
	"github.com/sixojke/test-astral/internal/service"
	"github.com/sixojke/test-astral/pkg/auth"
)

type Handler struct {
	service      *service.Service
	config       *config.Config
	tokenManager auth.TokenManager
}

func NewHandler(service *service.Service, config *config.Config, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		service:      service,
		config:       config,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(router *gin.RouterGroup) {
	register := router.Group("/register")
	{
		register.POST("", h.registerUser)
	}

	auth := router.Group("/auth")
	{
		auth.POST("", h.authUser)
		auth.DELETE("/:token", h.deleteSession)
	}

	docs := router.Group("/docs", h.middlewareAuth)
	{
		docs.POST("", h.uploadDocument)
		docs.GET("", h.getDocuments)
		docs.GET("/:id", h.getDocument)
		docs.HEAD("/:id", h.checkDocument)
		docs.DELETE("/:id", h.deleteDocument)
	}
}
