package v1

import (
	"github.com/gin-gonic/gin"
	"my_blog_backend/internal/usecase"
)

type Handler struct {
	services *usecase.Services
}

func NewHandler(services *usecase.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		//...
	}
}
