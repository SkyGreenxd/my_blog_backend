package domain

import (
	"my_blog_backend/pkg/e"
	"strings"
	"time"
	"unicode/utf8"
)

type Category struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func (c *Category) Validate() error {
	if err := validateCategoryName(c.Name); err != nil {
		return err
	}

	return nil
}

func validateCategoryName(name string) error {
	name = strings.TrimSpace(name)
	length := utf8.RuneCountInString(name)

	if length < 2 {
		return e.ErrCategoryTooShort
	}

	if length > 128 {
		return e.ErrCategoryTooLong
	}

	return nil
}
