package usecase

import (
	"my_blog_backend/internal/domain"
	"time"
)

type CreateUserReq struct {
	Username string
	Email    string
	Password string
}

type LoginUserReq struct {
	Email    string
	Password string
}

type UserRes struct {
	Username string
	Email    string
	Role     domain.Role
}

type LoginUserRes struct {
	SessionID             string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	User                  UserRes
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
