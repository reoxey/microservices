package core

import "context"

type UserRepo interface {
	All(ctx context.Context) (Users, error)
	ByID(ctx context.Context, id int) (*User, error)
	Add(ctx context.Context, user *User) (int, error)
	Edit(ctx context.Context, user *User) error
	Authenticate(ctx context.Context, email, pass string) (int, bool, error)
}
