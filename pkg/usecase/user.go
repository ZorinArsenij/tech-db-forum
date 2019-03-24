package usecase

import (
	userDomain "github.com/ZorinArsenij/tech-db-forum/pkg/domain/user"
	"github.com/ZorinArsenij/tech-db-forum/pkg/usecase/repository"
)

func NewUserInteractor(repo repository.User) *UserInteractor {
	return &UserInteractor{
		repository: repo,
	}
}

type UserInteractor struct {
	repository repository.User
}

func (i *UserInteractor) GetUserByNickname(nickname string) (*userDomain.User, error) {
	return i.repository.GetUserByNickname(nickname)
}

func (i *UserInteractor) UpdateUser(data *userDomain.Update, nickname string) (*userDomain.User, error) {
	return i.repository.UpdateUser(data, nickname)
}

func (i *UserInteractor) CreateUser(data *userDomain.User) (*userDomain.Users, error) {
	return i.repository.CreateUser(data)
}
