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
}

func NewCategory(name string) *Category {
	return &Category{
		Name: name,
	}
}

func (c *Category) ChangeName(newName string) error {
	if c.Name == newName {
		return e.ErrCategoryIsExists
	}

	c.Name = newName
	return nil
}

//func (c *Category) Validate() error {
//	if err := ValidateCategoryName(c.Name); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func ValidateCategoryName(name string) error {
//	name = strings.TrimSpace(name)
//	length := utf8.RuneCountInString(name)
//
//	if length < 2 {
//		return e.ErrCategoryTooShort
//	}
//
//	if length > 128 {
//		return e.ErrCategoryTooLong
//	}
//
//	return nil
//}
