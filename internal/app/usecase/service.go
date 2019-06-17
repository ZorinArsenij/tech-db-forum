package usecase

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/service"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/usecase/repository"
)

func NewServiceInteractor(repo repository.Service) *ServiceInteractor {
	return &ServiceInteractor{
		repository: repo,
	}
}

type ServiceInteractor struct {
	repository repository.Service
}

func (i *ServiceInteractor) GetStatus() (*service.Status, error) {
	return i.repository.GetStatus()
}
func (i *ServiceInteractor) Clear() error {
	return i.repository.Clear()
}
