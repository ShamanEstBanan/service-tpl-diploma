package usecase

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"go.uber.org/zap"
	"log"
	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
	"time"
)

func (s *service) CreateUser(ctx context.Context, user domain.NewUser) (authToken string, err error) {
	if user.Login == "" {
		return "", errs.ErrLoginIsEmpty
	} else if user.Password == "" {
		return "", errs.ErrPasswordIsEmpty
	}
	userID, err := s.storage.CreateUser(ctx, user)
	if err != nil {
		s.lg.Error("Error register user:", zap.Error(err))
		return "", err
	}
	token, err := generateToken(userID)
	if err != nil {
		s.lg.Error("Error generating token :", zap.Error(err))
		return "", err
	}
	return token, nil
}

func (s *service) AuthUser(ctx context.Context, user domain.AuthUser) (authToken string, err error) {
	if user.Login == "" {
		return "", errs.ErrLoginIsEmpty
	} else if user.Password == "" {
		return "", errs.ErrPasswordIsEmpty
	}
	userID, err := s.storage.CheckUser(ctx, user)
	if err != nil {
		s.lg.Error("Error auth user:", zap.Error(err))
		return "", err
	}
	if userID == "" {
		return "", errs.ErrInvalidLoginOrPassword
	}
	token, err := generateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// TODO убрать в env
const signingKey = "your_signing_key"

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

func generateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt: jwt.At(time.Now()),
		},
		Username: userID,
	})
	stringToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return stringToken, nil
}

func ParseToken(accessToken string, signingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, nil
	}
	return "", errs.ErrInvalidAccessToken
}
