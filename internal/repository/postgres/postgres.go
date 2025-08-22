package postgres

import (
	"my_blog_backend/pkg/e"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgDatabase struct {
	Db *gorm.DB
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

	return &PgDatabase{Db: db}, nil
}

func (pg *PgDatabase) Close() error {
	sqlDb, err := pg.Db.DB()
	if err != nil {
		return e.Wrap("failed to close db", err)
	}

	if err := sqlDb.Close(); err != nil {
		return e.Wrap("failed to close db", err)
	}

	return nil
}
