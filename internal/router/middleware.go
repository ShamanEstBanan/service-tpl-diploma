package router

import (
	"errors"
	"log"
	"net/http"
	"service-tpl-diploma/internal/errs"
	"service-tpl-diploma/internal/usecase"
	"strings"
)

// TODO убрать в env
const signingKey = "your_signing_key"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// пробуем вытащить куку
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			log.Printf("No auth token in headers")
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			log.Printf("invalid auth header: %s", headerParts)
			return
		}
		if headerParts[0] != "Bearer" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			log.Printf("invalid auth header: %s", headerParts)
			return
		}
		userID, err := usecase.ParseToken(headerParts[1], []byte(signingKey))
		if errors.Is(err, errs.ErrInvalidAccessToken) {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			log.Printf("invalid token: %s", headerParts)
			return
		}
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			log.Printf("invalid token: %s", headerParts)
			return
		}

		r.Header.Set("userId", userID)
		next.ServeHTTP(w, r)
	})
}
