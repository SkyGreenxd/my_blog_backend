package hash

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"my_blog_backend/pkg/e"
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
	const op = "BcryptHashManager.HashPassword"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), manager.cost)
	if err != nil {
		return "", e.Wrap(op, err)
	}

	return string(hash), nil
}

func (manager *BcryptHashManager) Compare(password string, hash string) error {
	const op = "BcryptHashManager.Compare"
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return e.Wrap(op, e.ErrMismatchedHashAndPassword)
		}
		return e.Wrap(op, err)
	}

	return nil
}

//func HashToken(token string) string {
//	// 1. Создаем новый объект хешера SHA-256.
//	hasher := sha256.New()
//
//	// 2. Записываем в хешер байты нашего токена.
//	// Важно передавать байты, а не просто строку.
//	hasher.Write([]byte(token))
//
//	// 3. Вычисляем хеш. `Sum(nil)` возвращает результат в виде среза байт.
//	hashBytes := hasher.Sum(nil)
//
//	// 4. Кодируем срез байт в шестнадцатеричную строку.
//	// [32]byte -> "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"
//	hashString := hex.EncodeToString(hashBytes)
//
//	return hashString
//}
