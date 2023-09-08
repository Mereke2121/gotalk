package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/gotalk/models"
	"github.com/gotalk/pkg/repository"
	utils "github.com/gotalk/utils"
)

var salt = "39@#$rkf@dk!dwk$#"

type AuthService struct {
	repo  *repository.Repository
	users map[string]*models.User
}

func NewAuthService(repo *repository.Repository) Authorization {
	//return &AuthService{
	//	users: make(map[string]*models.User),
	//	repo:  repo,
	//}
	return &AuthService{
		repo:  repo,
		users: make(map[string]*models.User),
	}
}

func (s *AuthService) AddUser(user *models.User) error {
	hashedPassword := hashPassword(user.Password)

	user.Password = hashedPassword
	err := s.repo.AddUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Authenticate(user *models.Authentication) (string, error) {
	user.Password = hashPassword(user.Password)

	id, err := s.repo.GetUserId(user)
	if err != nil {
		return "", err
	}

	log.Println("user id: ", id)

	// create jwt token and return in
	var jwtFields []utils.JWTTokenField
	jwtFields = append(jwtFields, utils.JWTTokenField{
		Type:  utils.UserId,
		Value: id,
	})
	return utils.CreateToken(jwtFields)
}

func (s *AuthService) GetUserById(userId string) (*models.User, error) {
	return s.repo.GetUserById(userId)
}

func hashPassword(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	hash := h.Sum([]byte(salt))
	return hex.EncodeToString(hash)
}
