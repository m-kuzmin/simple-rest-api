package db

import "context"

// Querier is for all queries to all tables in the DB
type Querier interface {
	UserQuerier
}

// UserQuerier is for queries to the users table
type UserQuerier interface {
	CreateUsers(context.Context, []User) error
	SearchUsers(_ context.Context, name, phoneNumber, country, city string) ([]User, error)
}

type User struct {
	ID          int64
	Name        string
	PhoneNumber string
	Country     string
	City        string
}
