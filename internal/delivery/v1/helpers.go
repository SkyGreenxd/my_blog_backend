package v1

import (
	"errors"
	"log"
	"my_blog_backend/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorToHttpRes(err error, c *gin.Context) {
	log.Println(err)

	var code int
	var message string

	switch {
	case errors.Is(err, e.ErrCategoryNotFound):
		code = http.StatusNotFound
		message = "category not found"
	case errors.Is(err, e.ErrArticleNotFound):
		code = http.StatusNotFound
		message = "article not found"
	case errors.Is(err, e.ErrUserNotAuthor):
		code = http.StatusUnauthorized
		message = "you are not the author of the post"
	case errors.Is(err, e.ErrNoDataToUpdate):
		code = http.StatusUnprocessableEntity
		message = "no data to update"
	case errors.Is(err, e.ErrArticleNameIsExists):
		code = http.StatusUnprocessableEntity
		message = "the name of the article is not changed"
	case errors.Is(err, e.ErrArticleContentIsExists):
		code = http.StatusUnprocessableEntity
		message = "the content of the article is not changed"
	case errors.Is(err, e.ErrCategoryIsExists):
		code = http.StatusUnprocessableEntity
		message = "category with such name already exists"
	case errors.Is(err, e.ErrUserNotFound):
		code = http.StatusNotFound
		message = "user not found"
	case errors.Is(err, e.ErrUsernameIsExists) || errors.Is(err, e.ErrEmailIsExists):
		code = http.StatusUnprocessableEntity
		message = "username or email is exists"
	case errors.Is(err, e.ErrInvalidCredentials):
		code = http.StatusUnauthorized
		message = "invalid credentials"
	case errors.Is(err, e.ErrUsernameIsSame):
		code = http.StatusUnprocessableEntity
		message = "username is same"
	case errors.Is(err, e.ErrEmailIsSame):
		code = http.StatusUnprocessableEntity
		message = "email is same"
	case errors.Is(err, e.ErrUsernameIsForbidden):
		code = http.StatusForbidden
		message = "username is forbidden"
	case errors.Is(err, e.ErrPasswordIsSame):
		code = http.StatusUnprocessableEntity
		message = "password is same"
	case errors.Is(err, e.ErrRefreshTokenHashDuplicate):
		code = http.StatusInternalServerError
		message = "the refresh token hash is duplicated"
	case errors.Is(err, e.ErrUnauthorized):
		code = http.StatusUnauthorized
		message = "unauthorized"
	case errors.Is(err, e.ErrArticleCategoryIsExists):
		code = http.StatusUnprocessableEntity
		message = "the category of the article is not changed"
	case errors.Is(err, e.ErrUserAlreadyAdmin):
		code = http.StatusUnprocessableEntity
		message = "user is already admin"
	case errors.Is(err, e.ErrPermissionDenied):
		code = http.StatusForbidden
		message = "permission denied"
	case errors.Is(err, e.ErrCategoryInUse):
		code = http.StatusUnprocessableEntity
		message = "category is in use"
	case errors.Is(err, e.ErrArticleDataIsInvalid):
		code = http.StatusUnprocessableEntity
		message = "title or content of the article is invalid"
	case errors.Is(err, e.ErrCategorySlugIsExists):
		code = http.StatusUnprocessableEntity
		message = "category slug is exists"
	default:
		code = http.StatusInternalServerError
		message = "internal server error"
	}

	c.JSON(code, gin.H{"error": message})
}
