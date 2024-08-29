package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/AsakoKabe/go-yandex-shortener/internal/app/context"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils/jwt"
)

// CookieName Ключ под которым хранится кука
const CookieName = "jwt"

// Auth Middleware для аутентификации пользователя и получения userID
func Auth(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var tokenString string

			cookie, err := r.Cookie(CookieName)

			if err != nil {
				tokenString, err = jwt.BuildJWTString(uuid.NewString())

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				http.SetCookie(
					w, &http.Cookie{
						Name:  CookieName,
						Value: tokenString,
					},
				)
			} else {
				tokenString = cookie.Value
			}

			userID, err := jwt.GetUserID(tokenString)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.SetUserID(r.Context(), userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}(next)
}
