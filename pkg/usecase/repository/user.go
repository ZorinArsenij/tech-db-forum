package repository

import (
	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/user"
)

type User interface {
	GetUserByNickname(nickname string) (*user.User, error)
	UpdateUser(data *user.Update, nickname string) (*user.User, error)
	CreateUser(data *user.User) (*user.Users, error)
}
