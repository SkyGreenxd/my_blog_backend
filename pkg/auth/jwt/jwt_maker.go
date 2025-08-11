package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"my_blog_backend/internal/usecase"
	"strconv"
	"time"
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

func (manager *TokenManager) Generate(id uint, email, role string) (string, error) {
	claims, err := NewUserClaims(id, email, role, manager.duration)
	if err != nil {
		log.Printf("error generating jwt: %v", err)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(manager.secretKey))
	if err != nil {
		log.Printf("error signing jwt: %v", err)
		return "", fmt.Errorf("error signing jwt: %w", err)
	}

	log.Printf("generated jwt: %s", tokenString)
	return tokenString, nil
}

func (manager *TokenManager) Verify(tokenString string) (*usecase.AuthPrincipal, error) {
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
		log.Printf("error parsing jwt: %v", err)
		return nil, fmt.Errorf("error parsing jwt: %w", err)
	}

	if !token.Valid {
		log.Printf("invalid jwt: %v", err)
		return nil, fmt.Errorf("invalid token")
	}

	log.Printf("parsed jwt: %s", tokenString)
	return claimsToAuthPrincipal(claims)
}

func claimsToAuthPrincipal(claims *UserClaims) (*usecase.AuthPrincipal, error) {
	userId, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		log.Printf("error parsing claims: %v", err)
		return nil, fmt.Errorf("invalid user id in token subject: %w", err)
	}

	authPrincipal := &usecase.AuthPrincipal{
		ID:    uint(userId),
		Email: claims.Email,
		Role:  claims.Role,
	}

	log.Printf("parsed auth principal: %v", authPrincipal)
	return authPrincipal, nil
}
