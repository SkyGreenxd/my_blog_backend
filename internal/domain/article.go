package domain

import (
	"time"
)

type Article struct {
	ID         uint
	Title      string
	Content    string
	AuthorID   uint
	CategoryID uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Author     *User
	Category   *Category
}

func (*Article) TableName() string {
	return "articles"
}

func NewArticle(title, content string, authorId, CategoryId uint) *Article {
	return &Article{
		Title:      title,
		Content:    content,
		AuthorID:   authorId,
		CategoryID: CategoryId,
	}
}

// TODO: Перенести это в validators
//func (a *Article) Validate() error {
//	if err := ValidateTitle(a.Title); err != nil {
//		return err
//	}
//
//	if err := ValidateContent(a.Content); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func ValidateTitle(title string) error {
//	title = strings.TrimSpace(title)
//	length := utf8.RuneCountInString(title)
//
//	if length < 3 {
//		return e.ErrTitleTooShort
//	}
//
//	if length > 100 {
//		return e.ErrTitleTooLong
//	}
//
//	if strings.Contains(title, "<") || strings.Contains(title, ">") {
//		return e.ErrTitleHasHTML
//	}
//
//	return nil
//}
//
//func ValidateContent(content string) error {
//	content = strings.TrimSpace(content)
//	length := utf8.RuneCountInString(content)
//
//	if length < 10 {
//		return e.ErrContentTooShort
//	}
//
//	if length > 16000 {
//		return e.ErrContentTooLong
//	}
//
//	if strings.Contains(content, "<script") {
//		return e.ErrContentHasScript
//	}
//	return nil
//}
