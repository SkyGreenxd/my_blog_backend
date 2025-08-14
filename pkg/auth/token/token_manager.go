package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"my_blog_backend/internal/domain"
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

func (manager *TokenManager) NewJWT(userID uint, email string, role domain.Role) (string, error) {
	claims, err := NewUserClaims(userID, email, role, manager.duration)
	if err != nil {
		log.Printf("error generating token: %v", err)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(manager.secretKey))
	if err != nil {
		log.Printf("error signing token: %v", err)
		return "", fmt.Errorf("error signing token: %w", err)
	}

	log.Printf("generated token: %s", tokenString)
	return tokenString, nil
}

func (manager *TokenManager) VerifyJWT(tokenString string) (*usecase.AuthPrincipal, error) {
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
		log.Printf("error parsing token: %v", err)
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		log.Printf("invalid token: %v", err)
		return nil, fmt.Errorf("invalid token")
	}

	log.Printf("parsed token: %s", tokenString)
	return claimsToAuthPrincipal(claims)
}

func (manager *TokenManager) NewRefreshToken() (string, string, error) {
	b := make([]byte, 32) // 256 бит
	_, err := rand.Read(b)
	if err != nil {
		return "", "", err
	}

	token := base64.URLEncoding.EncodeToString(b)
	hashed := manager.HashRefreshToken(token)
	return token, hashed, nil
}

func (manager *TokenManager) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
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
