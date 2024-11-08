package middleware

import (
	"context"
	"github.com/Fyefhqdishka/web-project/pkg/jwt"
	"log/slog"
	"net/http"
)

func JWTMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				logger.Warn("Токен не найден", "error", err)
				http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
				return
			}

			tokenStr := cookie.Value
			claims, err := jwt.VerifyJWT(tokenStr)
			if err != nil {
				logger.Warn("Недействительный токен", "error", err)
				http.Error(w, "Недействительный токен", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.ID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
