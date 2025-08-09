package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"my_blog_backend/pkg/e"
	"os"
)

type PgDatabase struct {
	db *gorm.DB
}

func Connect() (*PgDatabase, error) {
	path := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(path), &gorm.Config{})
	if err != nil {
		return nil, e.Wrap("failed to connect to db", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, e.Wrap("failed to get sql db instance", err)
	}

	if err := sqlDb.Ping(); err != nil {
		return nil, e.Wrap("failed to ping db", err)
	}

	return &PgDatabase{db: db}, nil
}
