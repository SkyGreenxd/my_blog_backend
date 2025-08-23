package domain

import (
	"my_blog_backend/pkg/e"
	"strings"
	"time"
)

type User struct {
	ID           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Role         Role
	Username     string
	Email        string
	PasswordHash string
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func NewUser(username, email, passwordHash string) *User {
	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         RoleUser,
	}
}

func (u *User) Validate() error {
	forbiddenUsernames := []string{"admin", "root", "user"}
	for _, forbidden := range forbiddenUsernames {
		if strings.EqualFold(u.Username, forbidden) {
			return e.ErrUsernameIsForbidden
		}
	}

	return nil
}

func (u *User) ChangePassword(newPasswordHash string) error {
	if u.PasswordHash == newPasswordHash {
		return e.ErrPasswordIsSame
	}
	u.PasswordHash = newPasswordHash
	return nil
}

func (u *User) ChangeUsername(newUsername string) error {
	if u.Username == newUsername {
		return e.ErrUsernameIsSame
	}

	u.Username = newUsername
	return nil
}

func (u *User) ChangeEmail(newEmail string) error {
	if u.Email == newEmail {
		return e.ErrEmailIsSame
	}

	u.Email = newEmail
	return nil
}

func (u *User) SetAdminRole() error {
	if u.Role == RoleAdmin {
		return e.ErrUserAlreadyAdmin
	}
	u.Role = RoleAdmin
	return nil
}
