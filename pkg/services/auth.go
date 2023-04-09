package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gotalk/models"
	"github.com/gotalk/utils"
	"github.com/pkg/errors"
	"log"
)

var users = make(map[string]models.User)
var salt = "39@#$rkf@dk!dwk$#"

func (s *Service) AddUser(user *models.User) error {
	hashedPassword := hashPassword(user.Password)

	if _, ok := users[user.Email]; !ok {
		users[user.Email] = models.User{
			UserName: user.UserName,
			Email:    user.Email,
			Password: hashedPassword,
		}
	} else {
		return errors.New("you already authorized")
	}

	log.Println("users: ", users)

	return nil
}

func (s *Service) Authenticate(user *models.Authentication) (string, error) {
	hashedPassword := hashPassword(user.Password)

	if _, ok := users[user.Email]; ok {
		if users[user.Email].Password != hashedPassword {
			return "", errors.New("unauthorized")
		}
	} else {
		return "", errors.Errorf("there's no user by this email: %s", user.Email)
	}

	log.Println("users: ", users)

	// create jwt token and return in
	var jwtFields []utils.JWTTokenField
	jwtFields = append(jwtFields, utils.JWTTokenField{
		Type:  utils.UserEmail,
		Value: user.Email,
	})
	return utils.CreateToken(jwtFields)
}

func hashPassword(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	hash := h.Sum([]byte(salt))
	return hex.EncodeToString(hash)
}
