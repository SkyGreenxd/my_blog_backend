package usecase

import (
	"my_blog_backend/internal/domain"
	"time"
)

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRes struct {
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     domain.Role `json:"role"`
}

type LoginUserRes struct {
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  UserRes   `json:"user"`
}

func (u *LoginUserReq) Validate() error {
	if err := domain.ValidateEmail(u.Email); err != nil {
		return err
	}

	if err := domain.ValidateUserPassword(u.Password); err != nil {
		return err
	}

	return nil
}

func (u *CreateUserReq) Validate() error {
	if err := domain.ValidateUsername(u.Username); err != nil {
		return err
	}

	if err := domain.ValidateEmail(u.Email); err != nil {
		return err
	}

	if err := domain.ValidateUserPassword(u.Password); err != nil {
		return err
	}

	return nil
}
