package core

import (
	"context"
)

type UserService interface {
	AllUsers(ctx context.Context) (Users, error)
	UserById(ctx context.Context, id int) (*User, error)
	AddUser(ctx context.Context, user *User) (int, error)
	EditUser(ctx context.Context, user *User) error
	Login(ctx context.Context, login *Login) (string, error)
	Authorize(s string) (map[string]interface{}, error)
}
