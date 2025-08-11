package domain

import (
	"my_blog_backend/pkg/e"
	"strings"
	"time"
	"unicode/utf8"
)

type Article struct {
	ID         uint
	Title      string
	Content    string
	AuthorID   uint
	CategoryID uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (*Article) TableName() string {
	return "articles"
}

func (a *Article) Validate() error {
	if err := validateTitle(a.Title); err != nil {
		return err
	}

	if err := validateContent(a.Content); err != nil {
		return err
	}

	return nil
}

func validateTitle(title string) error {
	title = strings.TrimSpace(title)
	length := utf8.RuneCountInString(title)

	if length < 3 {
		return e.ErrTitleTooShort
	}

	if length > 100 {
		return e.ErrTitleTooLong
	}

	if strings.Contains(title, "<") || strings.Contains(title, ">") {
		return e.ErrTitleHasHTML
	}

	return nil
}

func validateContent(content string) error {
	content = strings.TrimSpace(content)
	length := utf8.RuneCountInString(content)

	if length < 10 {
		return e.ErrContentTooShort
	}

	if length > 16000 {
		return e.ErrContentTooLong
	}

	if strings.Contains(content, "<script") {
		return e.ErrContentHasScript
	}
	return nil
}
