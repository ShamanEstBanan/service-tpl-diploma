package storage

import (
	"context"
	"service-tpl-diploma/internal/domain"
)

func (s *storage) GetAccountBalance(ctx context.Context, user domain.NewUser) error {
	query := "INSERT INTO users (login,password) VALUES($1,crypt($2,gen_salt('bf',8)))"

	_, err := s.db.Exec(ctx, query, user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}
