package usecase

import (
	"my_blog_backend/internal/domain"
	"time"
)

type AuthPrincipal struct {
	ID    uint
	Role  domain.Role
	Email string
}

type TokenResponse struct {
	Token     string
	ExpiresAt time.Time
}

type HashManager interface {
	HashPassword(password string) (string, error)
	Compare(password string, hash string) error
}

type TokenManager interface {
	NewJWT(userID uint, email string, role domain.Role) (*TokenResponse, error)
	VerifyJWT(tokenString string) (*AuthPrincipal, error)
	NewRefreshToken() (token string, hashed string, err error)
	HashRefreshToken(token string) string
}
