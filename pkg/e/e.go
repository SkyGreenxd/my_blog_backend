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
	ErrUsernameDuplicate    = errors.New("the user's username address already exists.")
	ErrEmailDuplicate       = errors.New("the user's email address already exists.")
	// role
	ErrInvalidRole = errors.New("invalid role")
	// email
	ErrEmailTooLong       = errors.New("email is too long")
	ErrEmailInvalidFormat = errors.New("email format is invalid")
	ErrEmailHasSpaces     = errors.New("email contains spaces")
	ErrEmailTooShort      = errors.New("email is too short")
	// password
	ErrPasswordTooShort  = errors.New("password is too short")
	ErrPasswordTooLong   = errors.New("password is too long")
	ErrPasswordHasSpaces = errors.New("password contains spaces")

	// categories
	ErrCategoryDuplicate = errors.New("category with such name already exists")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrCategoryTooShort  = errors.New("category name is empty")
	ErrCategoryTooLong   = errors.New("category name is too long")
	ErrCategoryInUse     = errors.New("category is already in use")

	// articles
	ErrTitleTooShort    = errors.New("title is too short")
	ErrTitleTooLong     = errors.New("title is too long")
	ErrTitleHasHTML     = errors.New("title has html")
	ErrContentTooShort  = errors.New("content is too short")
	ErrContentTooLong   = errors.New("content is too long")
	ErrContentHasScript = errors.New("content has script")
	ErrArticleNotFound  = errors.New("article not found")
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
