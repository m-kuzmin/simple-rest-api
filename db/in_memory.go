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

// CreateUsers implements Querier.
func (db *InMemoryDB) CreateUsers(_ context.Context, users []User) error {
	db.Users = append(db.Users, users...)

	logging.Debugf("InMemoryDB.CreateUsers called with args: %v", users)
	logging.Tracef("InMemoryBD.Users: %v", db.Users)

	return nil
}
