package access

import (
	"github.com/nynrathod/automator-api/pkg/entities"
)

type Service interface {
	ToggleAccess(data *entities.Access) (any, error)
	SharedUser(data *entities.Access) (any, error)
	SharedAccess(data *entities.Access) (any, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) ToggleAccess(data *entities.Access) (any, error) {

	//return nil, nil
	return s.repository.ToggleAccess(data)
}

func (s *service) SharedUser(data *entities.Access) (any, error) {

	return s.repository.SharedUser(data)
}

func (s *service) SharedAccess(data *entities.Access) (any, error) {

	return s.repository.SharedAccess(data)
}
