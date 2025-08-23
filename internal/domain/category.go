package domain

import (
	"my_blog_backend/pkg/e"
	"time"
)

type Category struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Slug      string
}

func NewCategory(name, slug string) *Category {
	return &Category{
		Name: name,
		Slug: slug,
	}
}

func (c *Category) ChangeName(newName string) error {
	if c.Name == newName {
		return e.ErrCategoryIsExists
	}

	c.Name = newName
	return nil
}

func (c *Category) ChangeSlug(newSlug string) error {
	if c.Slug == newSlug {
		return e.ErrCategorySlugIsExists
	}

	c.Slug = newSlug
	return nil
}
