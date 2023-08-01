package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/m-kuzmin/simple-rest-api/db/sqlc"
	"github.com/m-kuzmin/simple-rest-api/logging"
)

type Postgres struct {
	sqlDB *sql.DB
	conn  *sqlc.Queries
}

func NewPostgres(conn *sql.DB) *Postgres {
	return &Postgres{
		conn:  sqlc.New(conn),
		sqlDB: conn,
	}
}

// CreateUsers implements UserQuerier.
func (db *Postgres) CreateUsers(ctx context.Context, users []User) error {
	for _, user := range users {
		arg := sqlc.CreateUserParams{
			ID:          user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Country:     user.Country,
			City:        user.City,
		}

		if err := db.conn.CreateUser(ctx, arg); err != nil {
			return fmt.Errorf("PostgreSQL error: %w", err)
		}
	}

	return nil
}

func (db *Postgres) Close() error {
	if err := db.sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close sql connection handle: %w", err)
	}

	return nil
}

func PostgresMigrateUp(db *sql.DB, migrationsSource, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "migrations",
		DatabaseName:    dbName,
	})
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL database driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(migrationsSource, dbName, driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate client: %w", err)
	}

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate the database: %w", err)
	}

	return nil
}

func ConnectToDBWithRetry(postgresDriver, postgresAddr string, retries uint, interval time.Duration,
) (*sql.DB, error) {
	conn, err := sql.Open(postgresDriver, postgresAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create sql database connection to %q with driver %q: %w", postgresAddr,
			postgresDriver, err)
	}

	for i := uint(0); i < retries; i++ {
		err = conn.Ping()
		if err != nil {
			if i > 0 {
				logging.Warnf("Retry %d: Pinging PostgreSQL after %s because it has not started yet", i,
					interval.String())
			}

			time.Sleep(interval)

			continue
		}

		return conn, nil
	}

	return nil, fmt.Errorf("failed to ping database %q after %d attempts: %w", postgresAddr, retries, err)
}
