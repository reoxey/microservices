package mock

import (
	"context"

	"user/profile"
)

type repo struct {

}

func (r repo) All(ctx context.Context) (profile.Users, error) {
	return []profile.User{
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

func (r repo) ByID(ctx context.Context, i int) (profile.User, error) {
	return profile.User{
		Id:       i,
		Name:     "One",
		Email:    "one@test.com",
		IsAdmin:  true,
		JoinedAt: "Today",
	}, nil
}

func (r repo) Add(ctx context.Context, user *profile.User) (int, error) {
	return 0, nil
}

func (r repo) Edit(ctx context.Context, user *profile.User) error {
	return nil
}

func (r repo) Authenticate(ctx context.Context, email, pass string) (int, bool, error) {
	return 0, true, nil
}

func NewMock() profile.UserRepo {
	return &repo{}
}
