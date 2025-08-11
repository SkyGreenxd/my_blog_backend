package repository

import (
	"context"
	"my_blog_backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetById(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
}

type ArticleRepository interface {
	Create(ctx context.Context, article *domain.Article) (*domain.Article, error)
	GetByID(ctx context.Context, id uint) (*domain.Article, error)
	Update(ctx context.Context, article *domain.Article) error
	Delete(ctx context.Context, id uint) error
	ListAll(ctx context.Context) ([]domain.Article, error)
	ListByAuthor(ctx context.Context, authorID uint) ([]domain.Article, error)
	ListByCategory(ctx context.Context, categoryID uint) ([]domain.Article, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) (*domain.Category, error)
	GetByID(ctx context.Context, id uint) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uint) error
	ListAll(ctx context.Context) ([]domain.Category, error)
}
