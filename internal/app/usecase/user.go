package usecase

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/user"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/usecase/repository"
)

func NewUserInteractor(repo repository.User) *UserInteractor {
	return &UserInteractor{
		repository: repo,
	}
}

type UserInteractor struct {
	repository repository.User
}

func (i *UserInteractor) GetUserByNickname(nickname string) (*user.User, error) {
	return i.repository.GetUserByNickname(nickname)
}

func (i *UserInteractor) UpdateUser(data *user.Update, nickname string) (*user.User, error) {
	return i.repository.UpdateUser(data, nickname)
}

func (i *UserInteractor) CreateUser(data *user.User) (*user.Users, error) {
	return i.repository.CreateUser(data)
}
