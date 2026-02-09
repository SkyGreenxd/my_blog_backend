package v1

import (
	"log"
	"my_blog_backend/internal/delivery"
	"my_blog_backend/pkg/e"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary		Sign Up
// @Tags			auth
// @Description	Create a new user account
// @Accept			json
// @Produce		json
// @Param			input	body		delivery.CreateUserRequest	true	"Sign up data"
// @Success		201		{object}	delivery.UserRes
// @Failure		400		{object}	delivery.ErrResponse
// @Failure		500		{object}	delivery.ErrResponse
// @Router			/auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var req delivery.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	user, err := h.services.UserService.CreateUser(c.Request.Context(), delivery.ToServiceCreateUserReq(&req))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusCreated, delivery.ToUserRes(user))
}

// @Summary		Sign In
// @Tags			auth
// @Description	Log in to an existing account
// @Accept			json
// @Produce		json
// @Param			input	body		delivery.LoginRequest	true	"Login data"
// @Success		200		{object}	delivery.LoginUserRes
// @Failure		400		{object}	delivery.ErrResponse
// @Failure		401		{object}	delivery.ErrResponse
// @Router			/auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var req delivery.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	res, err := h.services.UserService.LoginUser(c.Request.Context(), delivery.ToLoginUserReq(&req))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, delivery.ToLoginUserRes(res))
}

// @Summary		Get Current User
// @Tags			users
// @Description	Get data of the currently authenticated user
// @Accept			json
// @Produce		json
// @Success		200	{object}	delivery.UserRes
// @Failure		401	{object}	delivery.ErrResponse
// @Failure		500	{object}	delivery.ErrResponse
// @Security		ApiKeyAuth
// @Router			/users/me [get]
func (h *Handler) getCurrentUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			ErrorToHttpRes(e.ErrUnauthorized, c)
		} else {
			ErrorToHttpRes(e.ErrInternalServer, c)
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

func (h *Handler) getUserById(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	user, err := h.services.UserService.GetUserById(c.Request.Context(), uint(userId))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, delivery.ToUserRes(user))
}

// @Summary		Get User By Username
// @Tags			users
// @Description	Find user by their username
// @Accept			json
// @Produce		json
// @Param			username	path		string	true	"Username"
// @Success		200			{object}	delivery.UserRes
// @Failure		404			{object}	delivery.ErrResponse
// @Router			/users/{username} [get]
func (h *Handler) getUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.services.UserService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, delivery.ToUserRes(user))
}

// @Summary		Update User
// @Tags			users
// @Description	Update current user's profile information
// @Accept			json
// @Produce		json
// @Param			input	body		delivery.UpdateUserReq	true	"New user data"
// @Success		200		{object}	delivery.UserRes
// @Failure		400		{object}	delivery.ErrResponse
// @Failure		401		{object}	delivery.ErrResponse
// @Security		ApiKeyAuth
// @Router			/users/me/update [patch]
func (h *Handler) updateUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			ErrorToHttpRes(e.ErrUnauthorized, c)
		} else {
			ErrorToHttpRes(e.ErrInternalServer, c)
		}
		return
	}

	var newData delivery.UpdateUserReq
	if err := c.ShouldBindJSON(&newData); err != nil {
		log.Println(err)
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	newUser, err := h.services.UserService.UpdateUser(c.Request.Context(), userId.(uint), delivery.ToUpdateUserReq(&newData))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, delivery.ToUserRes(newUser))
}

// @Summary		Change Password
// @Tags			auth
// @Description	Change current user's password
// @Accept			json
// @Produce		json
// @Param			input	body		delivery.ChangePasswordReq	true	"Password change data"
// @Success		200		{object}	map[string]interface{}
// @Failure		400		{object}	delivery.ErrResponse
// @Failure		401		{object}	delivery.ErrResponse
// @Security		ApiKeyAuth
// @Router			/auth/password/change [post]
func (h *Handler) changePassword(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			ErrorToHttpRes(e.ErrUnauthorized, c)
		} else {
			ErrorToHttpRes(e.ErrInternalServer, c)
		}
		return
	}

	var req delivery.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	if err := h.services.UserService.ChangePassword(c.Request.Context(), userId.(uint), delivery.ToChangePasswordReq(&req)); err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary		Refresh Session
// @Tags			auth
// @Description	Obtain a new access token using a refresh token
// @Accept			json
// @Produce		json
// @Param			input	body		delivery.RefreshTokenReq	true	"Refresh token"
// @Success		200		{object}	delivery.LoginUserRes
// @Failure		400		{object}	delivery.ErrResponse
// @Failure		401		{object}	delivery.ErrResponse
// @Router			/auth/refresh [post]
func (h *Handler) refreshSession(c *gin.Context) {
	var req delivery.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	res, err := h.services.UserService.RefreshSession(c.Request.Context(), req.RefreshToken)
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, delivery.ToLoginUserRes(res))
}

// @Summary		Logout
// @Tags			auth
// @Description	Invalidate current session
// @Accept			json
// @Produce		json
// @Param			input	body		delivery.LogoutUserReq	true	"Refresh token to invalidate"
// @Success		200		{object}	map[string]interface{}
// @Failure		400		{object}	delivery.ErrResponse
// @Router			/auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	var req delivery.LogoutUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ErrorToHttpRes(e.ErrBadRequest, c)
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
			ErrorToHttpRes(e.ErrUnauthorized, c)
		} else {
			ErrorToHttpRes(e.ErrInternalServer, c)
		}
		return
	}

	if err := h.services.UserService.SetAdminRole(c.Request.Context(), userId.(uint)); err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
