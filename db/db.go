package db

import "context"

type Querier interface {
	UserQuerier
}

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
