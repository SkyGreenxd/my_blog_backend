package hash

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHashManager struct {
	cost int
}

func NewBcryptHashManager(cost int) (*BcryptHashManager, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("invalid bcrypt cost: %d, must be between %d and %d", cost, bcrypt.MinCost, bcrypt.MaxCost)
	}
	return &BcryptHashManager{cost: cost}, nil
}

func (manager *BcryptHashManager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), manager.cost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	return string(hash), nil
}

func (manager *BcryptHashManager) Compare(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
