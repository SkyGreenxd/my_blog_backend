package postgres

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"my_blog_backend/pkg/e"
)

func checkChangeQueryResult(result *gorm.DB, op, message string, notFound error) error {
	if err := result.Error; err != nil {
		return e.WrapDBError(op, err)
	}

	if result.RowsAffected == 0 {
		return notFound
	}

	log.Printf("%s: %s", op, message)
	return nil
}

func checkGetQueryResult(result *gorm.DB, op, message string, notFound error) error {
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return notFound
	}

	if err := result.Error; err != nil {
		return e.WrapDBError(op, err)
	}

	log.Printf("%s: %s", op, message)
	return nil
}
