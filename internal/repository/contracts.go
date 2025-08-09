package repository

import (
	"context"
	"my_blog_backend/internal/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetById(ctx context.Context, id uint) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	//Update(ctx context.Context, user *entities.User) error
	//Delete(ctx context.Context, id uint) error
}

type ArticleRepository interface {
	Create(ctx context.Context, article *entities.Article) error
	GetByID(ctx context.Context, id uint) (*entities.Article, error)
	Update(ctx context.Context, article *entities.Article) error
	Delete(ctx context.Context, id uint) error
	ListAll(ctx context.Context) ([]entities.Article, error)
	ListByAuthor(ctx context.Context, authorID uint) ([]entities.Article, error)
	ListByCategory(ctx context.Context, categoryID uint) ([]entities.Article, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category *entities.Category) error
	GetByID(ctx context.Context, id uint) (*entities.Category, error)
	Update(ctx context.Context, category *entities.Category) error
	Delete(ctx context.Context, id uint) error
	ListAll(ctx context.Context) ([]entities.Category, error)
}
