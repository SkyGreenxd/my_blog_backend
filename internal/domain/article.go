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
}
