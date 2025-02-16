package domain

import (
	"context"
)

type User struct {
	ID   string
	Name string
}

type (
	QueryOptions struct {
		Limit  *uint32
		Offset *uint32
	}
)

type UserRepository interface {
	Get(ctx context.Context, id string) (*User, error)
	Search(ctx context.Context, query QueryOptions) ([]*User, error)
	Create(ctx context.Context, user User) (*string, error)
	Delete(ctx context.Context, id string) error
}
