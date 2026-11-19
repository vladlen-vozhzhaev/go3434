package main

import (
	"context"
	"net/http"
)

type contextKey string

const userIDKey contextKey = "userID"

func userIdFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("session")
		if err != nil {
			http.Redirect(writer, request, "/login", http.StatusFound)
			return
		}
		userID, ok := getUserIDByToken(cookie.Value)
		if !ok {
			http.SetCookie(writer, &http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				HttpOnly: false,
				MaxAge:   -1, // сутки
			})
			http.Redirect(writer, request, "/login", http.StatusFound)
			return
		}
		ctx := context.WithValue(request.Context(), userIDKey, userID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	}

}
