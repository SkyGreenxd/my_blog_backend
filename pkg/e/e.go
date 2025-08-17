package e

import (
	"errors"
	"fmt"
)

var (
	// users
	ErrUserNotFound        = errors.New("user not found")
	ErrUserDuplicate       = errors.New("user with such email or username already exists")
	ErrUsernameIsForbidden = errors.New("username is forbidden")
	ErrPasswordIsSame      = errors.New("password is same")
	ErrUserAlreadyAdmin    = errors.New("user is already admin")
	// username
	ErrUsernameEmpty        = errors.New("username is empty")
	ErrUsernameTooShort     = errors.New("username is too short")
	ErrUsernameTooLong      = errors.New("username is too long")
	ErrUsernameInvalidChars = errors.New("username contains invalid characters")
	ErrUsernameHasSpaces    = errors.New("username contains spaces")
	ErrUsernameIsExists     = errors.New("the user's username already exists.")
	ErrEmailIsExists        = errors.New("the user's email address already exists.")
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
	ErrCategoryIsExists = errors.New("category with such name already exists")
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryTooShort = errors.New("category name is empty")
	ErrCategoryTooLong  = errors.New("category name is too long")
	ErrCategoryInUse    = errors.New("category is already in use")

	// articles
	ErrTitleTooShort    = errors.New("title is too short")
	ErrTitleTooLong     = errors.New("title is too long")
	ErrTitleHasHTML     = errors.New("title has html")
	ErrContentTooShort  = errors.New("content is too short")
	ErrContentTooLong   = errors.New("content is too long")
	ErrContentHasScript = errors.New("content has script")
	ErrArticleNotFound  = errors.New("article not found")

	ErrMismatchedHashAndPassword = errors.New("password does not match hash")

	// Sessions
	ErrSessionRevoked            = errors.New("session revoked")
	ErrSessionExpired            = errors.New("session expired")
	ErrRefreshTokenHashDuplicate = errors.New("refresh token hash already exists")
	ErrSessionNotFound           = errors.New("session not found")
	ErrRefreshTokenInvalid       = errors.New("refresh token is invalid")

	// Общие ошибки
	ErrInvalidEmail       = errors.New("invalid email")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrParseFailed        = errors.New("parse failed")
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
