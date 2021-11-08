package http

import (
	"log"
	"net/http"

	v1 "github.com/Feokrat/user-balance-api/internal/delivery/http/v1"

	"github.com/Feokrat/user-balance-api/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
	logger   *log.Logger
}

func NewHandler(services *service.Services, logger *log.Logger) *Handler {
	return &Handler{services: services, logger: logger}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.logger)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
