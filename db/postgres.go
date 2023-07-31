package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/m-kuzmin/simple-rest-api/db/sqlc"
)

type Postgres struct {
	conn *sqlc.Queries
}

func NewPostgres(dbSource string) (*Postgres, error) {
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		return nil, fmt.Errorf("while connecting to postgres (%q): %w", dbSource, err)
	}

	return &Postgres{
		conn: sqlc.New(conn),
	}, nil
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
