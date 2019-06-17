package postgresql

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/thread"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/vote"
	"github.com/jackc/pgx"
)

const (
	createVote = "createVote"
	getVote    = "getVote"
	updateVote = "updateVote"
)

var voteQueries = map[string]string{
	createVote: `INSERT INTO vote (voice, user_nickname, thread_id)
	VALUES ($1, $2, $3)`,

	getVote: `SELECT id, voice
	FROM vote
	WHERE user_nickname = $1 AND thread_id = $2`,

	updateVote: `UPDATE vote
	SET voice = $1
	WHERE id = $2`,
}

func NewVoteRepo(conn *pgx.ConnPool) *Vote {
	return &Vote{
		conn: conn,
	}
}

type Vote struct {
	conn *pgx.ConnPool
}

func (v *Vote) CreateVote(data *vote.Vote, slugOrId string) (*thread.Thread, error) {
	tx, err := v.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var threadID, voteID uint64
	var currentVote bool

	if err := tx.QueryRow(getUserByNickname, data.UserNickname).
		Scan(&data.UserNickname); err != nil {
		return nil, err
	}

	var slug *string
	if err := tx.QueryRow(checkThreadByIdOrSlug, slugOrId).
		Scan(&threadID, &slug); err != nil {
		return nil, err
	}

	if err := tx.QueryRow(getVote, data.UserNickname, threadID).Scan(&voteID, &currentVote); err == nil {
		if currentVote != data.Voice {
			if _, err := tx.Exec(updateVote, data.Voice, voteID); err != nil {
				return nil, err
			}
			data.Rating *= 2
		} else {
			data.Rating = 0
		}
	} else {
		if _, err = tx.Exec(createVote, data.Voice, data.UserNickname, threadID); err != nil {
			return nil, err
		}

	}

	var received thread.Thread

	if err := tx.QueryRow(updateThreadVotes, data.Rating, threadID).
		Scan(&received.ID, &received.Slug, &received.Title, &received.Message, &received.ForumSlug, &received.UserNickname, &received.Created, &received.Votes); err != nil {
		return nil, err
	}

	tx.Commit()
	return &received, nil
}
