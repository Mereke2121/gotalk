package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const secretKey = "dsljiowdm#@DJ!da"

type JWTTokenField struct {
	Type  HeaderParam
	Value interface{}
}

// CreateToken creates jwt token
func CreateToken(tokenFields []JWTTokenField) (string, error) {
	// Устанавливаем время истечения токена
	expirationTime := time.Now().Add(time.Hour)

	// initialize jwt fields
	jwtMap := jwt.MapClaims{}
	jwtMap["exp"] = expirationTime.Unix()
	for _, field := range tokenFields {
		jwtMap[string(field.Type)] = field.Value
	}

	// create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtMap)

	// signing with our secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Функция проверки JWT токена
func VerifyToken(tokenString string, paramType HeaderParam) (interface{}, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// Возвращаем секретный ключ для проверки подписи
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}

	// Получаем идентификатор пользователя из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, jwt.ErrInvalidKey
	}

	result, ok := claims[string(paramType)]
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	return result, nil
}
