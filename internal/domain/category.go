package domain

import (
	"my_blog_backend/pkg/e"
	"strings"
	"time"
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
	if strings.TrimSpace(name) == "" {
		return e.ErrCategoryNameEmpty
	}

	if len(name) > 128 {
		return e.ErrCategoryNameLong
	}

	return nil
}
