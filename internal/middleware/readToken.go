package middleware

import (
	"context"
	"log"
	"net/http"
)

// 从http header里面获取token并放到context
func ReadToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		username := r.Header.Get("Username")
		token := r.Header.Get("Token")

		if username != "" && token != "" {
			ctx = context.WithValue(ctx, "username", username)
			ctx = context.WithValue(ctx, "token", token)
			log.Println("token", username, token)
		}

		// next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
