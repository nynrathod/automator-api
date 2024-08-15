package users

import (
	"github.com/nynrathod/automator-api/pkg/entities"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	VerifyEmail(email string) (*entities.User, error)
	Register(user *entities.User) (*entities.User, error)
	GetUser(email string) (*entities.User, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) VerifyEmail(email string) (*entities.User, error) {
	return s.repository.VerifyEmail(email)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *service) Register(user *entities.User) (*entities.User, error) {
	//encryptedUser, err := encryptUserData(user)
	//if err != nil {
	//	return encryptedUser, err
	//}

	hash, _ := hashPassword(user.Password)

	user.Password = hash

	return s.repository.Register(user)
}

func (s *service) GetUser(email string) (*entities.User, error) {
	return s.repository.GetUser(email)
}
