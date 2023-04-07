package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const secretKey = "dsljiowdm#@DJ!da"

// Функция создания JWT токена
func CreateToken(roomId int) (string, error) {
	// Устанавливаем время истечения токена
	expirationTime := time.Now().Add(time.Hour)

	// Создаем новый токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"roomId": roomId,
		"exp":    expirationTime.Unix(),
	})

	// Подписываем токен нашим секретным ключом
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Функция проверки JWT токена
func VerifyToken(tokenString string) (int, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// Возвращаем секретный ключ для проверки подписи
		return secretKey, nil
	})
	if err != nil {
		return 0, err
	}

	// Получаем идентификатор пользователя из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, jwt.ErrInvalidKey
	}

	roomId, ok := claims["roomId"].(float64)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	return int(roomId), nil
}
