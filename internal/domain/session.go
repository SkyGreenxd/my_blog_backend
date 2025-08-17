package domain

import (
	"my_blog_backend/pkg/e"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id               uuid.UUID
	UserId           uint
	RefreshTokenHash string
	IsRevoked        bool
	CreatedAt        time.Time
	ExpiresAt        time.Time
}

func NewSession(userID uint, refreshTokenHash string, expiresAt time.Time) *Session {
	return &Session{
		Id:               uuid.New(),
		UserId:           userID,
		RefreshTokenHash: refreshTokenHash,
		IsRevoked:        false,
		CreatedAt:        time.Now().UTC(),
		ExpiresAt:        expiresAt,
	}
}

func (s *Session) ValidateState() error {
	// Отозвана ли сессия?
	if s.IsRevoked {
		return e.ErrSessionRevoked
	}

	// Истек ли срок действия?
	if s.ExpiresAt.Before(time.Now().UTC()) {
		return e.ErrSessionExpired
	}

	return nil
}
