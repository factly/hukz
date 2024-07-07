package util

import (
	"context"
	"errors"
	"net/http"
)

type ctxKeyUserID string

var UserIDKey ctxKeyUserID

// CheckUser checks X-User in header
func CheckUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User")
		if user == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser returns user ID as string
func GetUser(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", errors.New("context not found")
	}
	userID := ctx.Value(UserIDKey)

	if userID != nil {
		return "", nil
	}
	return "", errors.New("something went wrong")
}
