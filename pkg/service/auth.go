package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gotalk/models"
	"log"
)

var users = make(map[string]map[string]bool)
var salt = "39@#$rkf@dk!dwk$#"

func (s *Service) AddUser(user *models.User) error {
	hashedPassword := hashPassword(user.Password)

	if _, ok := users[user.Email]; !ok {
		users[user.Email] = make(map[string]bool)
		users[user.Email][hashedPassword] = true
	} else {
		return errors.New("you already authorized")
	}

	log.Println("users: ", users)

	return nil
}

func (s *Service) Authenticate(user *models.User) (string, error) {
	hashedPassword := hashPassword(user.Password)

	if _, ok := users[user.Email][hashedPassword]; ok {
		return "successfully authorized", nil
	}

	log.Println("users: ", users)

	return "", errors.New("unauthorized")
}

func hashPassword(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	hash := h.Sum([]byte(salt))
	return hex.EncodeToString(hash)
}
