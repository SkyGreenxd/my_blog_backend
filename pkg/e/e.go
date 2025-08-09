package e

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserDublicate     = errors.New("user with such email or username already exists")
	ErrCategoryDublicate = errors.New("category with such name already exists")
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
