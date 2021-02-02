package mysql

import (
	"context"
	"fmt"
	"strings"

	"user/profile"
)

func (m mysqlRepo) All(ctx context.Context) (profile.Users, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT id, name, email, is_admin, joined_at FROM "+m.table+" WHERE 1")
	if err != nil {
		return nil, err
	}

	var users []profile.User
	for rows.Next() {
		var user profile.User

		if err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.IsAdmin, &user.JoinedAt); err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (m mysqlRepo) ByID(ctx context.Context, id int) (profile.User, error) {

	var user profile.User

	err := m.db.QueryRowContext(ctx, "SELECT id, name, email, is_admin, joined_at FROM "+m.table+" WHERE id=?", id).
		Scan(&user.Id, &user.Name, &user.Email, &user.IsAdmin, &user.JoinedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m mysqlRepo) Add(ctx context.Context, user *profile.User) (int, error) {

	s := strings.Builder{}
	s.WriteString("name = '"+user.Name+"'")
	s.WriteString(",email = '"+user.Email+"'")
	s.WriteString(",password = '"+user.Password+"'")
	is := "0"
	if user.IsAdmin {
		is = "1"
	}
	s.WriteString(",is_admin = '"+is+"'")

	rows, err := m.db.ExecContext(ctx, "INSERT "+m.table+" SET "+s.String())

	if err != nil {
		return 0, err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return 0, NoRowsAffected
	}

	id, _ := rows.LastInsertId()

	return int(id), nil
}

func (m mysqlRepo) Edit(ctx context.Context, user *profile.User) error {

	var s []string
	if user.Id == 0 {
		return noIdUpdate
	}
	if user.Name != "" {
		s = append(s, "name = '"+user.Name+"'")
	}
	if user.Email != "" {
		s = append(s, "email = '"+user.Email+"'")
	}
	if user.Password != "" {
		s = append(s, "password = '"+user.Password+"'")
	}

	rows, err := m.db.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", m.table, strings.Join(s, ","), user.Id),
	)

	if err != nil {
		return err
	}

	if n, _ := rows.RowsAffected(); n == 0 {
		return NoRowsAffected
	}

	return err
}

func (m mysqlRepo) Authenticate(ctx context.Context, email, pass string) (int, bool, error) {

	var id, is int

	err := m.db.QueryRowContext(ctx, "SELECT id, is_admin FROM "+m.table+" WHERE email=? AND password=?", email, pass).
		Scan(&id, &is)

	if err != nil {
		return 0, false, err
	}

	isAdmin := false
	if is == 1 {
		isAdmin = true
	}

	return id, isAdmin, nil
}
