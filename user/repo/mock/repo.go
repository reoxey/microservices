package mock

import (
	"context"

	"user/core"
)

type repo struct {
}

func (r repo) All(ctx context.Context) (core.Users, error) {
	return []core.User{
		{
			Id:       1,
			Name:     "One",
			Email:    "one@test.com",
			IsAdmin:  true,
			JoinedAt: "Today",
		},
		{
			Id:       2,
			Name:     "Two",
			Email:    "two@test.com",
			IsAdmin:  false,
			JoinedAt: "Today",
		},
	}, nil
}

func (r repo) ByID(ctx context.Context, i int) (core.User, error) {
	return core.User{
		Id:       i,
		Name:     "One",
		Email:    "one@test.com",
		IsAdmin:  true,
		JoinedAt: "Today",
	}, nil
}

func (r repo) Add(ctx context.Context, user *core.User) (int, error) {
	return 0, nil
}

func (r repo) Edit(ctx context.Context, user *core.User) error {
	return nil
}

func (r repo) Authenticate(ctx context.Context, email, pass string) (int, bool, error) {
	return 0, true, nil
}

func NewMock() core.UserRepo {
	return &repo{}
}
