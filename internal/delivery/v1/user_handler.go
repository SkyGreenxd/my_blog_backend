package v1

import (
	"log"
	"my_blog_backend/internal/delivery"
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
		ErrorToHttpRes(err, c)
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
		ErrorToHttpRes(err, c)
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
		ErrorToHttpRes(err, c)
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
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, delivery.ToUserRes(user))
}

func (h *Handler) getUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.services.UserService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		ErrorToHttpRes(err, c)
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
		ErrorToHttpRes(err, c)
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
		ErrorToHttpRes(err, c)
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
		ErrorToHttpRes(err, c)
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
		ErrorToHttpRes(err, c)
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
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
