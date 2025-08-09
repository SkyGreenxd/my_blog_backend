package e

import (
	"errors"
	"fmt"
	"log"
)

var (
	// users
	ErrUserNotFound  = errors.New("user not found")
	ErrUserDuplicate = errors.New("user with such email or username already exists")
	// username
	ErrUsernameEmpty        = errors.New("username is empty")
	ErrUsernameTooShort     = errors.New("username is too short")
	ErrUsernameTooLong      = errors.New("username is too long")
	ErrUsernameInvalidChars = errors.New("username contains invalid characters")
	ErrUsernameHasSpaces    = errors.New("username contains spaces")
	// role
	ErrInvalidRole = errors.New("invalid role")
	// email
	ErrEmailTooLong       = errors.New("email is too long")
	ErrEmailInvalidFormat = errors.New("email format is invalid")
	ErrEmailHasSpaces     = errors.New("email contains spaces")
	ErrEmailTooShort      = errors.New("email is too short")
	// password
	ErrPasswordTooShort  = errors.New("password hash is too short")
	ErrPasswordTooLong   = errors.New("password hash is too long")
	ErrPasswordHasSpaces = errors.New("password hash contains spaces")

	// categories
	ErrCategoryDuplicate = errors.New("category with such name already exists")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrCategoryNameEmpty = errors.New("category name is empty")
	ErrCategoryNameLong  = errors.New("category name is too long")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapDBError(op string, err error) error {
	if err == nil {
		return nil
	}
	wrapped := Wrap(op, err)
	log.Print(wrapped)
	return wrapped
}
