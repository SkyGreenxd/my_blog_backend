package v1

import (
	"errors"
	"my_blog_backend/internal/delivery"
	"my_blog_backend/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser(c *gin.Context) {
	var newUser delivery.CreateUserRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	UserRes, err := h.services.UserService.CreateUser(c.Request.Context(), delivery.ToServiceCreateUserReq(&newUser))
	if err != nil {
		switch {
		case errors.Is(err, e.ErrUsernameIsExists):
			c.JSON(http.StatusConflict, gin.H{"error": e.ErrUsernameIsExists.Error()})
		case errors.Is(err, e.ErrEmailIsExists):
			c.JSON(http.StatusConflict, gin.H{"error": e.ErrEmailIsExists.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": e.ErrInternalServer.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, delivery.ToUserRes(UserRes))
}

func (h *Handler) Login(c *gin.Context) {
	var authorizedUser delivery.LoginRequest
	if err := c.ShouldBindJSON(&authorizedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	res, err := h.services.UserService.LoginUser(c.Request.Context(), delivery.ToLoginUserReq(&authorizedUser))
	if err != nil {
		switch {
		case errors.Is(err, e.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": e.ErrInvalidCredentials.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": e.ErrInternalServer.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, delivery.ToLoginUserRes(res))
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		}
		return
	}

	user, err := h.services.UserService.GetUserById(c.Request.Context(), userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, delivery.ToUserRes(user))
}
