package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"my_blog_backend/internal/domain"
	"strconv"
	"time"
)

type UserClaims struct {
	Email string      `json:"email"`
	Role  domain.Role `json:"role"`
	jwt.RegisteredClaims
}

func NewUserClaims(userId uint, email string, role domain.Role, duration time.Duration) (*UserClaims, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token ID: %w", err)
	}

	return &UserClaims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Subject:   strconv.FormatUint(uint64(userId), 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
