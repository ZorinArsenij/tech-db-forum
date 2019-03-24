package usecase

import (
	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/thread"
	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/vote"
	"github.com/ZorinArsenij/tech-db-forum/pkg/usecase/repository"
)

func NewVoteInteractor(repo repository.Vote) *VoteInteractor {
	return &VoteInteractor{
		repository: repo,
	}
}

type VoteInteractor struct {
	repository repository.Vote
}

func (i *VoteInteractor) CreateVote(data *vote.Vote, slugOrId string) (*thread.Thread, error) {
	return i.repository.CreateVote(data, slugOrId)
}
