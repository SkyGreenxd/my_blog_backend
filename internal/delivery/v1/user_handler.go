package v1

import (
	"errors"
	"log"
	"my_blog_backend/internal/delivery"
	"my_blog_backend/pkg/e"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TODO: убрать повторяющийся код
func (h *Handler) signUp(c *gin.Context) {
	var req delivery.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	user, err := h.services.UserService.CreateUser(c.Request.Context(), delivery.ToServiceCreateUserReq(&req))
	if err != nil {
		log.Println(err)
		if errors.Is(err, e.ErrUsernameIsExists) || errors.Is(err, e.ErrEmailIsExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "username or email is exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, delivery.ToUserRes(user))
}

// TODO: возвращать данные в куки
func (h *Handler) signIn(c *gin.Context) {
	var req delivery.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	res, err := h.services.UserService.LoginUser(c.Request.Context(), delivery.ToLoginUserReq(&req))
	if err != nil {
		log.Println(err)
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

func (h *Handler) getCurrentUser(c *gin.Context) {
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
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, delivery.ToUserRes(user))
}

// TODO: должны еще выводиться в будущем посты пользователя
func (h *Handler) getUserById(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.services.UserService.GetUserById(c.Request.Context(), uint(userId))
	if err != nil {
		log.Println(err)
		if errors.Is(err, e.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, delivery.ToUserRes(user))
}

func (h *Handler) updateUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		}
		return
	}

	var newData delivery.UpdateUserReq
	if err := c.ShouldBindJSON(&newData); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	newUser, err := h.services.UserService.UpdateUser(c.Request.Context(), userId.(uint), delivery.ToUpdateUserReq(&newData))
	if err != nil {
		log.Println(err)
		if errors.Is(err, e.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if errors.Is(err, e.ErrNoDataToUpdate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no data to update"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, delivery.ToUserRes(newUser))
}

func (h *Handler) changePassword(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		}
		return
	}

	var req delivery.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := h.services.UserService.ChangePassword(c.Request.Context(), userId.(uint), delivery.ToChangePasswordReq(&req)); err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, e.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, e.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": e.ErrInvalidCredentials.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// TODO: возвращать данные в куки
func (h *Handler) refreshSession(c *gin.Context) {
	var req delivery.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	res, err := h.services.UserService.RefreshSession(c.Request.Context(), req.RefreshToken)
	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, e.ErrUnauthorized):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, delivery.ToLoginUserRes(res))
}

func (h *Handler) logout(c *gin.Context) {
	var req delivery.LogoutUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := h.services.UserService.LogoutUser(c.Request.Context(), req.RefreshToken); err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, e.ErrUnauthorized):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) setAdminRole(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		}
		return
	}

	if err := h.services.UserService.SetAdminRole(c.Request.Context(), userId.(uint)); err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, e.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, e.ErrUserAlreadyAdmin):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is already admin"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}
}
