package v1

import (
	"log"

	"github.com/Feokrat/user-balance-api/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
	logger   *log.Logger
}

func NewHandler(services *service.Services, logger *log.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUserBalanceRoutes(v1)
	}
}
