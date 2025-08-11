package domain

import (
	"my_blog_backend/pkg/e"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
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

func (*User) TableName() string {
	return "users"
}

func (u *User) Validate() error {
	if err := validateUsername(u.Username); err != nil {
		return err
	}

	if err := validateEmail(u.Email); err != nil {
		return err
	}

	if u.Role != RoleAdmin && u.Role != RoleUser {
		return e.ErrInvalidRole
	}

	return nil
}

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	length := utf8.RuneCountInString(username)

	if username == "" {
		return e.ErrUsernameEmpty
	}

	if length < 5 {
		return e.ErrUsernameTooShort
	}

	if length > 32 {
		return e.ErrUsernameTooLong
	}

	if strings.Contains(username, " ") {
		return e.ErrUsernameHasSpaces
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`).MatchString(username) {
		return e.ErrUsernameInvalidChars
	}

	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	length := utf8.RuneCountInString(email)

	if length < 3 {
		return e.ErrEmailTooShort
	}

	if length > 320 {
		return e.ErrEmailTooLong
	}

	if strings.Contains(email, " ") {
		return e.ErrEmailHasSpaces
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return e.ErrEmailInvalidFormat
	}

	return nil
}

func ValidateUserPassword(password string) error {
	if err := validatePassword(password); err != nil {
		return err
	}

	return nil
}

func validatePassword(password string) error {
	password = strings.TrimSpace(password)

	if len(password) < 8 {
		return e.ErrPasswordTooShort
	}

	if len(password) > 128 {
		return e.ErrPasswordTooLong
	}

	if strings.Contains(password, " ") {
		return e.ErrPasswordHasSpaces
	}

	return nil
}

// Доп проверка пароля
//var (
//	digit    = regexp.MustCompile(`[0-9]`)
//	lower    = regexp.MustCompile(`[a-z]`)
//	upper    = regexp.MustCompile(`[A-Z]`)
//	special  = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':",\.<>\/\?\\|]`)
//)
//
//if !digit.MatchString(password) {
//return e.ErrPasswordNoDigit
//}
//if !lower.MatchString(password) {
//return e.ErrPasswordNoLower
//}
//if !upper.MatchString(password) {
//return e.ErrPasswordNoUpper
//}
//if !special.MatchString(password) {
//return e.ErrPasswordNoSpecial
//}
