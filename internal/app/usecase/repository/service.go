package repository

import "github.com/ZorinArsenij/tech-db-forum/internal/app/domain/service"

type Service interface {
	GetStatus() (*service.Status, error)
	Clear() error
}
