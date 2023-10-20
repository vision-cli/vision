package auth

import (
	"context"
	"net/http"
)

type User struct {
	ID       string
	Username string
	Email    string
}

type contextKey string

var contextKeyUser = contextKey("authed-user")

func Authenicate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Authenticate user
		u := &User{
			ID:       "123",
			Username: "johndoe",
			Email:    "john@doe.com",
		}
		r = r.WithContext(WithUser(r.Context(), u))
		next.ServeHTTP(w, r)
	})
}

func UserFromContext(ctx context.Context) *User {
	return ctx.Value(contextKeyUser).(*User)
}

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, contextKeyUser, user)
}
