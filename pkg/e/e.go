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
	ErrUsernameIsSame      = errors.New("username is same")
	ErrUserAlreadyAdmin    = errors.New("user is already admin")
	ErrEmailIsSame         = errors.New("email is same")
	// username
	ErrUsernameInvalidChars = errors.New("username contains invalid characters")
	ErrUsernameHasSpaces    = errors.New("username contains spaces")
	ErrUsernameIsExists     = errors.New("the user's username already exists.")
	ErrEmailIsExists        = errors.New("the user's email address already exists.")
	// role
	ErrInvalidRole = errors.New("invalid role")
	// email
	ErrEmailInvalidFormat = errors.New("email format is invalid")
	ErrEmailHasSpaces     = errors.New("email contains spaces")
	// password
	ErrPasswordHasSpaces = errors.New("password contains spaces")

	// categories
	ErrCategoryIsExists     = errors.New("category with such name already exists")
	ErrCategorySlugIsExists = errors.New("category slug already exists")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryInUse        = errors.New("category is already in use")

	// articles
	ErrTitleHasHTML            = errors.New("title has html")
	ErrContentHasScript        = errors.New("content has script")
	ErrArticleNotFound         = errors.New("article not found")
	ErrArticleNameIsExists     = errors.New("the name of the article is not changed")
	ErrArticleContentIsExists  = errors.New("the content of the article is not changed")
	ErrArticleCategoryIsExists = errors.New("the category of the article is not changed")
	ErrArticleDataIsInvalid    = errors.New("title or content of the article is invalid")
	ErrUserNotAuthor           = errors.New("user is not author")
	ErrArticleDuplicate        = errors.New("article is duplicate")

	ErrMismatchedHashAndPassword = errors.New("password does not match hash")

	// Sessions
	ErrSessionRevoked            = errors.New("session revoked")
	ErrSessionExpired            = errors.New("session expired")
	ErrRefreshTokenHashDuplicate = errors.New("refresh token hash already exists")
	ErrSessionNotFound           = errors.New("session not found")
	ErrRefreshTokenInvalid       = errors.New("refresh token is invalid")

	// Общие ошибки
	ErrPermissionDenied   = errors.New("permission denied")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrParseFailed        = errors.New("parse failed")
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoDataToUpdate     = errors.New("no data to update")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
