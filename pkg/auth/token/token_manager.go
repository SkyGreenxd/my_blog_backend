package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"my_blog_backend/internal/domain"
	"my_blog_backend/internal/usecase"
	"my_blog_backend/pkg/e"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secretKey string
	duration  time.Duration
}

func NewTokenManager(secretKey string, duration time.Duration) *TokenManager {
	return &TokenManager{
		secretKey: secretKey,
		duration:  duration,
	}
}

func (manager *TokenManager) NewJWT(userID uint, email string, role domain.Role) (*usecase.TokenResponse, error) {
	const op = "TokenManager.NewJWT"
	expiresAt := time.Now().Add(manager.duration)

	claims, err := NewUserClaims(userID, email, role, expiresAt)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(manager.secretKey))
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return &usecase.TokenResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}, nil
}

func (manager *TokenManager) VerifyJWT(tokenString string) (*usecase.AuthenticatedUser, error) {
	const op = "tokenManager.VerifyJWT"
	claims := &UserClaims{}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(manager.secretKey), nil
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		keyFunc, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, e.Wrap(op, e.ErrParseFailed)
	}

	if !token.Valid {
		return nil, e.Wrap(op, e.ErrTokenInvalid)
	}

	return claimsToAuthPrincipal(claims)
}

func (manager *TokenManager) NewRefreshToken() (string, string, error) {
	const op = "tokenManager.NewRefreshToken"
	b := make([]byte, 32) // 256 бит
	_, err := rand.Read(b)
	if err != nil {
		return "", "", e.Wrap(op, err)
	}

	token := base64.URLEncoding.EncodeToString(b)
	hashed := manager.HashRefreshToken(token)
	return token, hashed, nil
}

func (manager *TokenManager) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func claimsToAuthPrincipal(claims *UserClaims) (*usecase.AuthenticatedUser, error) {
	const op = "tokenManager.claimsToAuthPrincipal"
	userId, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	authPrincipal := &usecase.AuthenticatedUser{
		ID:    uint(userId),
		Email: claims.Email,
		Role:  claims.Role,
	}

	return authPrincipal, nil
}
