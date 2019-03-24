package usecase

import (
	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/thread"
	"github.com/ZorinArsenij/tech-db-forum/pkg/usecase/repository"
)

func NewThreadInteractor(repo repository.Thread) *ThreadInteractor {
	return &ThreadInteractor{
		repository: repo,
	}
}

type ThreadInteractor struct {
	repository repository.Thread
}

func (i *ThreadInteractor) GetThread(slugOrId string) (*thread.Thread, error) {
	return i.repository.GetThread(slugOrId)
}

func (i *ThreadInteractor) GetThreads(slug string, limit *int, since *string, orderDesc bool) (*thread.Threads, error) {
	return i.repository.GetThreads(slug, limit, since, orderDesc)
}

func (i *ThreadInteractor) CreateThread(data *thread.Create) (*thread.Thread, error) {
	return i.repository.CreateThread(data)
}

func (i *ThreadInteractor) UpdateThread(data *thread.Update, slugOrId string) (*thread.Thread, error) {
	return i.repository.UpdateThread(data, slugOrId)
}
