package db

import (
	"context"

	"github.com/m-kuzmin/simple-rest-api/logging"
)

const preallocateUsers = 100

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Users: make([]User, 0, preallocateUsers),
	}
}

type InMemoryDB struct {
	Users []User
}

// CreateUsers implements UserQuerier.
func (db *InMemoryDB) CreateUsers(_ context.Context, users []User) error {
	db.Users = append(db.Users, users...)

	logging.Debugf("InMemoryDB.Users: %v", db.Users)

	return nil
}

func (*InMemoryDB) SearchUsers(_ context.Context, _, _, _, _ string) ([]User, error) {
	panic("Use the real database instead of the mock one")
}
