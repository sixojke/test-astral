package delivery

import (
	"net/http"

	_ "github.com/sixojke/test-astral/docs"

	"github.com/gin-gonic/gin"
	"github.com/sixojke/test-astral/internal/config"
	v1 "github.com/sixojke/test-astral/internal/delivery/v1"
	"github.com/sixojke/test-astral/internal/service"
	"github.com/sixojke/test-astral/pkg/auth"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func (h *Handler) Init() *gin.Engine {
	// Create a new router
	router := gin.Default()

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", h.ping)

	h.initAPI(router)

	return router
}

// Test route to check server functionality
func (h *Handler) ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// Initialize api of several versions
func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.service, h.config, h.tokenManager)
	api := router.Group("/api")
	handlerV1.Init(api)
}
