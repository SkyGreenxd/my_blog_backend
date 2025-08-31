package v1

import (
	"errors"
	"my_blog_backend/internal/usecase"
	"my_blog_backend/pkg/e"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	tokenManager usecase.TokenManager
}

func NewMiddleware(tokenManager usecase.TokenManager) *Middleware {
	return &Middleware{tokenManager: tokenManager}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": e.ErrUnauthorized.Error(),
			})
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}
		jwtToken := strings.TrimPrefix(authHeader, prefix)

		authenticatedUser, err := m.tokenManager.VerifyJWT(jwtToken)
		if err != nil {
			if errors.Is(err, e.ErrTokenInvalid) || errors.Is(err, e.ErrParseFailed) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": e.ErrUnauthorized.Error(),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": e.ErrInternalServer.Error(),
			})
			return
		}

		c.Set("user_id", authenticatedUser.ID)
		c.Set("role", authenticatedUser.Role)

		c.Next()
	}
}
