package core

import (
	"context"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type userService struct {
	repo     UserRepo
	cache	 Cache
	auth	 JWTService
	validate *validator.Validate
}

func (u userService) Authorize(s string) (map[string]interface{}, error) {
	return u.auth.ValidateToken(s)
}

func (u userService) AllUsers(ctx context.Context) (Users, error) {
	var users Users
	err := u.cache.GetJSON(ctx, "all_users", &users)
	if err != nil {
		users, err = u.repo.All(ctx)
		if err != nil {
			return nil, err
		}
		u.cache.SetJSON(ctx, "all_users", &users, 0)
	}
	return users, nil
}

func (u userService) UserById(ctx context.Context, id int) (*User, error) {
	var user *User
	err := u.cache.GetJSON(ctx, "user_"+strconv.Itoa(id), &user)
	if err != nil {
		return u.repo.ByID(ctx, id)
	}
	return user, nil
}

func (u userService) AddUser(ctx context.Context, user *User) (id int, err error) {

	if err = u.validate.Struct(user); err != nil {
		return
	}
	if id, err = u.repo.Add(ctx, user); err != nil {
		return
	}
	u.cache.SetJSON(ctx,  "user_"+strconv.Itoa(id), &user, 0)
	return
}

func (u userService) EditUser(ctx context.Context, user *User) error {
	err := u.repo.Edit(ctx, user)
	if err != nil {
		return err
	}
	return u.cache.SetJSON(ctx,  "user_"+strconv.Itoa(user.Id), &user, 0)
}

func (u userService) Login(ctx context.Context, login *Login) (string, error) {
	id, isAdmin, err := u.repo.Authenticate(ctx, login.Email, login.Password)
	if err != nil {
		return "", err
	}
	return u.auth.GenerateToken(id, login.Email, isAdmin)
}

func NewService(ur UserRepo, cache Cache, auth JWTService) UserService {
	return &userService{
		repo: ur,
		cache: cache,
		auth: auth,
		validate: validator.New(),
	}
}
