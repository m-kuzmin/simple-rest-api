package db

import (
	"context"
	"strings"

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

func (db *InMemoryDB) SearchUsers(_ context.Context, name, phoneNumber, country, city string) ([]User, error) {
	results := make([]User, 0, len(db.Users))

	logging.Debugf("SearchUsers args: name=%q phoneNumber=%q country=%q city=%q", name, phoneNumber, country, city)

	for _, user := range db.Users {
		if name != "" && !strings.Contains(user.Name, name) {
			continue
		}

		if phoneNumber != "" && !strings.Contains(user.PhoneNumber, phoneNumber) {
			continue
		}

		if country != "" && !strings.Contains(user.Country, country) {
			continue
		}

		if city != "" && !strings.Contains(user.City, city) {
			continue
		}

		logging.Debugf("SearchUsers selected: %v", user)

		results = append(results, user)
	}

	return results, nil
}
