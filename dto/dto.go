package dto

import "my_blog_backend/internal/domain"

type CreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *CreateUser) Validate() error {
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
