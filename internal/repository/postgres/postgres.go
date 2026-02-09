package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"my_blog_backend/pkg/e"
	"os"

	"github.com/golang-migrate/migrate/v4"
	mpgdb "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgDatabase struct {
	Db  *gorm.DB
	dsn string
}

func Connect() (*PgDatabase, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

	return &PgDatabase{Db: db, dsn: dsn}, nil
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

func (db *PgDatabase) RunMigrations() error {
	const op = "PgDatabase.RunMigrations"

	sqlDb, err := sql.Open("pgx", db.dsn)
	if err != nil {
		return err
	}
	defer sqlDb.Close()

	driver, err := mpgdb.WithInstance(sqlDb, &mpgdb.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return e.Wrap(op, err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return e.Wrap(op, err)
	}

	log.Println("migrations applied successfully")
	return nil
}
