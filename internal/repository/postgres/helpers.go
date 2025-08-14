package postgres

import (
	"errors"
	"gorm.io/gorm"
	"my_blog_backend/pkg/e"
)

func checkChangeQueryResult(result *gorm.DB, op string, notFound error) error {
	if err := result.Error; err != nil {
		return e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return e.Wrap(op, notFound)
	}

	return nil
}

func checkGetQueryResult(result *gorm.DB, notFound error) error {
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return notFound
	}

	if err := result.Error; err != nil {
		return err
	}

	return nil
}
