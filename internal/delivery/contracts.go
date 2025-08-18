package delivery

import (
	"my_blog_backend/internal/domain"
	"my_blog_backend/internal/usecase"
	"time"
)

// TODO: добавить полную валидацию, чтобы сервис получал точно корректно-введенные данные
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=5,max=32,nospaces"`
	Email    string `json:"email" binding:"required,email,min=3,max=320,nospaces"`
	Password string `json:"password" binding:"required,min=8,max=128,nospaces"`
}

type UserRes struct {
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     domain.Role `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,min=3,max=320,nospaces"`
	Password string `json:"password" binding:"required,min=8,max=128,nospaces"`
}

type ErrResponse struct {
	ErrCode string `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type LoginUserRes struct {
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  UserRes   `json:"user"`
}

type UpdateUserReq struct {
	Username *string `json:"username" binding:"required,min=5,max=32,nospaces"`
	Email    *string `json:"email" binding:"required,email,min=3,max=32,nospaces"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=128,nospaces"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128,nospaces"`
}

func ToChangePasswordReq(req *ChangePasswordReq) *usecase.ChangePasswordReq {
	return &usecase.ChangePasswordReq{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
}

func ToUpdateUserReq(req *UpdateUserReq) *usecase.UpdateUserReq {
	return &usecase.UpdateUserReq{
		Username: req.Username,
		Email:    req.Email,
	}
}

func ToLoginUserRes(res *usecase.LoginUserRes) *LoginUserRes {
	return &LoginUserRes{
		SessionID:             res.SessionID,
		AccessToken:           res.AccessToken,
		RefreshToken:          res.RefreshToken,
		AccessTokenExpiresAt:  res.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: res.RefreshTokenExpiresAt,
		User:                  *ToUserRes(&res.User),
	}
}

func ToServiceCreateUserReq(request *CreateUserRequest) *usecase.CreateUserReq {
	return &usecase.CreateUserReq{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
}

func ToUserRes(res *usecase.UserRes) *UserRes {
	return &UserRes{
		Username: res.Username,
		Email:    res.Email,
		Role:     res.Role,
	}
}

func ToLoginUserReq(req *LoginRequest) *usecase.LoginUserReq {
	return &usecase.LoginUserReq{
		Email:    req.Email,
		Password: req.Password,
	}
}
