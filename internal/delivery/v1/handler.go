package v1

import (
	"github.com/gin-gonic/gin"
	"my_blog_backend/internal/usecase"
	"my_blog_backend/pkg/auth/token"
)

type Handler struct {
	services     *usecase.Services
	tokenManager token.TokenManager
}

func NewHandler(services *usecase.Services, tokenManager token.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		//...
	}
}
