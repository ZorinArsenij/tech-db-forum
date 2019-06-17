package usecase

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/forum"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/user"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/usecase/repository"
)

func NewForumInteractor(repo repository.Forum) *ForumInteractor {
	return &ForumInteractor{
		repository: repo,
	}
}

type ForumInteractor struct {
	repository repository.Forum
}

func (i *ForumInteractor) GetForum(slug string) (*forum.Forum, error) {
	return i.repository.GetForum(slug)
}

func (i *ForumInteractor) CreateForum(data *forum.Create) (*forum.Forum, error) {
	return i.repository.CreateForum(data)
}

func (i *ForumInteractor) GetForumUsers(slug string, limit *int, since *string, orderDesc bool) (*user.Users, error) {
	return i.repository.GetForumUsers(slug, limit, since, orderDesc)
}
