package repository

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/forum"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/user"
)

type Forum interface {
	GetForum(slug string) (*forum.Forum, error)
	CreateForum(data *forum.Create) (*forum.Forum, error)
	GetForumUsers(slug string, limit *int, since *string, orderDesc bool) (*user.Users, error)
}
