package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"service-tpl-diploma/internal/errs"

	"service-tpl-diploma/internal/domain"
)

func (s *storage) CreateUser(ctx context.Context, user domain.NewUser) (accountID string, err error) {
	// TODO: Сделать через транзакции
	query := "INSERT INTO users (login,password) VALUES($1,crypt($2,gen_salt('bf',8)))"
	_, err = s.db.Exec(ctx, query, user.Login, user.Password)
	if err != nil {
		return "", err
	}
	switch e := err.(type) {
	case *pgconn.PgError:
		switch e.Code {
		case pgerrcode.UniqueViolation:
			return "", errs.ErrLoginAlreadyExist
		}
	}
	query = fmt.Sprintf("SELECT id FROM users WHERE login ='%s'", user.Login)
	err = s.db.QueryRow(ctx, query).Scan(&accountID)
	if err != nil {
		return "", err
	}

	query = "INSERT INTO accounts (id) VALUES($1)"
	_, err = s.db.Exec(ctx, query, accountID)
	if err != nil {
		return accountID, err
	}
	return accountID, nil
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
