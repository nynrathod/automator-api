package users

import (
	"fmt"
	"github.com/nynrathod/automator-api/pkg/entities"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	VerifyEmail(email string, isLogin bool) (*entities.User, error)
	Register(user *entities.User) (*entities.User, error)
	GetUser(email string) (*entities.User, error)
	AddUser(data *entities.User) (*entities.User, error)

	ListUser(data string) (*entities.User, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) VerifyEmail(email string, isLogin bool) (*entities.User, error) {
	return s.repository.VerifyEmail(email, isLogin)
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

func (s *service) AddUser(data *entities.User) (*entities.User, error) {
	fmt.Println("servicedata", data)
	return s.repository.AddUser(data)
}

func (s *service) ListUser(data string) (*entities.User, error) {
	fmt.Println("servicedata", data)
	return s.repository.ListUser(data)
}
