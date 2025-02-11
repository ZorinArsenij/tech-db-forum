package repository

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/thread"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/vote"
)

type Vote interface {
	CreateVote(data *vote.Vote, slugOrId string) (*thread.Thread, error)
}
