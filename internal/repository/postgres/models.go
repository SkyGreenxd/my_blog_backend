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
	Username     string      `gorm:"size:32;uniqueIndex:idx_username;not null"`
	Email        string      `gorm:"size:320;uniqueIndex:idx_email;not null"`
	PasswordHash string      `gorm:"not null"`
}

type ArticleModel struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string `gorm:"size:128;not null"`
	Content    string `gorm:"not null"`
	AuthorID   uint   `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CategoryID uint   `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type CategoryModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:128;unique;not null"`
}

// TODO: реализовать хуки AfterDelete/BeforeDelete
