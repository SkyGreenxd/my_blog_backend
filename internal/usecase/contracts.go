package usecase

import "my_blog_backend/internal/domain"

type AuthPrincipal struct {
	ID    uint
	Role  domain.Role
	Email string
}

type HashManager interface {
	HashPassword(password string) (string, error)
	Compare(password string, hash string) error
}

type TokenManager interface {
	NewJWT(userID uint, email string, role domain.Role) (string, error)
	VerifyJWT(tokenString string) (*AuthPrincipal, error)
	NewRefreshToken() (token string, hashed string, err error)
	HashRefreshToken(token string) string
}
