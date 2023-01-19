package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"log"
	"service-tpl-diploma/internal/app/domain"
	"time"
)

var (
	ErrInvalidAccessToken = errors.New("invalid auth token")
)

func (s *service) CreateUser(ctx context.Context, user domain.NewUser) error {
	if user.Login == "" {
		return errors.New("login is empty")
	} else if user.Password == "" {
		return errors.New("password is empty")
	}
	err := s.storage.CreateUser(ctx, user)
	if err != nil {
		s.lg.Sugar().Error(err.Error())
		return err
	}
	return nil
}

func (s *service) AuthUser(ctx context.Context, user domain.AuthUser) (authToken string, err error) {
	if user.Login == "" {
		return "", errors.New("login is empty")
	} else if user.Password == "" {
		return "", errors.New("password is empty")
	}
	userId, err := s.storage.CheckUser(ctx, user)
	if err != nil {
		s.lg.Sugar().Error(err.Error())
		return "", err
	}
	if userId == "" {
		return "", errors.New("login or password is invalid")
	}
	token, err := generateToken(userId)
	if err != nil {
		return "", err
	}
	fmt.Println(token)

	return token, nil
}

// const hashSalt = "vsohAFzfiyAPFadu24n"
const signingKey = "your_signing_key"

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

func generateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Username: userId,
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
	return "", ErrInvalidAccessToken
}
