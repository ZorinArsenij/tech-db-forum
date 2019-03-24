package repository

import (
	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/thread"
	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/vote"
)

type Vote interface {
	CreateVote(data *vote.Vote, slugOrId string) (*thread.Thread, error)
}
