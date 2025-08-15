package usecase

import (
	"my_blog_backend/internal/domain"
)

type HashManager interface {
	HashPassword(password string) (string, error)
	Compare(password string, hash string) error
}

type TokenManager interface {
	NewJWT(userID uint, email string, role domain.Role) (*TokenResponse, error)
	VerifyJWT(tokenString string) (*AuthenticatedUser, error)
	NewRefreshToken() (token string, hashed string, err error)
	HashRefreshToken(token string) string
}
