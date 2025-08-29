package domain

import (
	"my_blog_backend/pkg/e"
	"strings"
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

func NewArticle(title, content string, authorId, CategoryId uint) *Article {
	return &Article{
		Title:      title,
		Content:    content,
		AuthorID:   authorId,
		CategoryID: CategoryId,
	}
}

func (a *Article) Validate() error {
	if err := ValidateTitle(a.Title); err != nil {
		return err
	}

	if err := ValidateContent(a.Content); err != nil {
		return err
	}

	return nil
}

func ValidateTitle(title string) error {
	if strings.Contains(title, "<") || strings.Contains(title, ">") {
		return e.ErrTitleHasHTML
	}

	return nil
}

func ValidateContent(content string) error {
	if strings.Contains(content, "<script") {
		return e.ErrContentHasScript
	}

	return nil
}

func (a *Article) ChangeTitle(newTitle string) error {
	if newTitle == a.Title {
		return e.ErrArticleNameIsExists
	}

	a.Title = newTitle
	return nil
}

func (a *Article) ChangeContent(newContent string) error {
	if newContent == a.Content {
		return e.ErrArticleContentIsExists
	}

	a.Content = newContent
	return nil
}

func (a *Article) ChangeCategory(newCategory *Category) error {
	if a.Category != nil && newCategory.ID == a.Category.ID {
		return e.ErrCategoryIsExists
	}

	a.Category = newCategory
	return nil
}

func (a *Article) CheckAuthor(userId uint) error {
	if a.AuthorID != userId {
		return e.ErrUserNotAuthor
	}

	return nil
}
