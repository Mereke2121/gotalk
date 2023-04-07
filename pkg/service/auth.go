package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gotalk/models"
)

var users []models.User
var salt = "39@#$rkf@dk)!dwk$#"

func (s *Service) AddUser(user *models.User) {
	hashedPassword := hashPassword(user.Password)
	user.Password = hashedPassword

	// add user
	users = append(users, *user)
}

func hashPassword(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	hash := h.Sum([]byte(salt))
	return hex.EncodeToString(hash)
}
