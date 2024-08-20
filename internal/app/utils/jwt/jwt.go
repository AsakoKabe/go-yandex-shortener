package jwt

import (
	"fmt"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Claims Данные для jwt
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// TokenExp Время жизни токена для пользователя
const TokenExp = time.Hour * 3

// SecretKey Секрет для создания JWT
var SecretKey = utils.GetEnv("JWT_TOKEN", "somevalue")

// BuildJWTString Создать JWT токен по id пользователя
func BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
			},
			UserID: userID,
		},
	)

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserID Получить id пользователя из JWT токена
func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("not valid token")
	}
	return claims.UserID, nil
}
