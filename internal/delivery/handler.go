package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/sixojke/test-astral/internal/delivery/v1"
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

func (h *Handler) Init() *gin.Engine {
	// Create a new router
	router := gin.Default()

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
	handlerV1 := v1.NewHandler(h.service)
	api := router.Group("/api")
	handlerV1.Init(api)
}
