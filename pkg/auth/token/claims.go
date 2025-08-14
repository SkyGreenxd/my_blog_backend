package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"my_blog_backend/internal/domain"
	"my_blog_backend/pkg/e"
	"strconv"
	"time"
)

type UserClaims struct {
	Email string      `json:"email"`
	Role  domain.Role `json:"role"`
	jwt.RegisteredClaims
}

func NewUserClaims(userId uint, email string, role domain.Role, expiresAt time.Time) (*UserClaims, error) {
	const op = "token.NewUserClaims"
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return &UserClaims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Subject:   strconv.FormatUint(uint64(userId), 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}, nil
}
