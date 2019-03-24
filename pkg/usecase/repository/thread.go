package repository

import "github.com/ZorinArsenij/tech-db-forum/pkg/domain/thread"

type Thread interface {
	GetThread(slugOrId string) (*thread.Thread, error)
	GetThreads(slug string, limit *int, since *string, orderDesc bool) (*thread.Threads, error)
	CreateThread(data *thread.Create) (*thread.Thread, error)
	UpdateThread(data *thread.Update, slugOrId string) (*thread.Thread, error)
}
