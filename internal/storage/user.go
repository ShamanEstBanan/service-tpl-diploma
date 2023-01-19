package storage

import (
	"context"
	"fmt"
	"service-tpl-diploma/internal/app/domain"
)

func (s *storage) CreateUser(ctx context.Context, user domain.NewUser) error {
	query := "INSERT INTO users (login,password) VALUES($1,crypt($2,gen_salt('bf',8)))"

	_, err := s.db.Exec(ctx, query, user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) CheckUser(ctx context.Context, user domain.AuthUser) (userID string, err error) {
	var exist bool
	query := fmt.Sprintf(
		"SELECT (password = crypt('%s', password)) AS pswmatch FROM users WHERE login = '%s'",
		user.Password, user.Login)

	err = s.db.QueryRow(ctx, query).Scan(&exist)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", err
	}
	query = fmt.Sprintf("SELECT id FROM users WHERE login = '%s'", user.Login)
	err = s.db.QueryRow(ctx, query).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, err
}

func (s *storage) GetUserId(ctx context.Context, login string) (id string, err error) {
	query := fmt.Sprintf("SELECT id FROM users WHERE login = '%s'", login)

	err = s.db.QueryRow(ctx, query).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
