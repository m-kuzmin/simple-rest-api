package db

import "context"

// Querier is for all queries to all tables in the DB
type Querier interface {
	UserQuerier
}

// UserQuerier is for queries to the users table
type UserQuerier interface {
	CreateUsers(context.Context, []User) error
}

type User struct {
	Name        string
	PhoneNumber string
	Country     string
	City        string
	ID          int64
}
