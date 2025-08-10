package postgres

import (
	"my_blog_backend/internal/domain"
	"time"
)

type UserModel struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Role         domain.Role `gorm:"not null"`
	Username     string      `gorm:"uniqueIndex:idx_username;not null"`
	Email        string      `gorm:"uniqueIndex:idx_email;not null"`
	PasswordHash string      `gorm:"not null"`
}

type ArticleModel struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string `gorm:"not null"`
	Content    string
	AuthorID   uint `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CategoryID uint `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type CategoryModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"unique;not null"`
}

// TODO: реализовать хуки AfterDelete/BeforeDelete
