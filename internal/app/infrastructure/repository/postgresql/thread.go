package postgresql

import (
	"errors"

	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/thread"

	"github.com/jackc/pgx"
)

const (
	getThreadBySlug = `SELECT id, slug, title, message, forum_slug, user_nickname, created, votes
	FROM thread
	WHERE slug = $1;`

	getThreadById = `SELECT id, slug, title, message, forum_slug, user_nickname, created, votes
	FROM thread
	WHERE id = $1;`

	getThreadByIdOrSlug = `SELECT id, slug, title, message, forum_slug, user_nickname, created, votes
	FROM thread
	WHERE slug = $1 OR id::TEXT = $1`

	getThreadShortBySlugOrId = `SELECT id, forum_slug
	FROM thread
	WHERE slug = $1 OR id::TEXT = $1`

	createThread = `INSERT INTO thread (slug, title, message, forum_id, forum_slug, user_id, user_nickname, created)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	)
	RETURNING id, slug, title, message, forum_slug, user_nickname, created, votes;`

	getThreadsByForumSlugLimit = `SELECT id, slug, title, message, forum_slug, user_nickname, created, votes
	FROM thread
	WHERE forum_slug = $1
	ORDER BY created
 	LIMIT $2;`

	getThreadsByForumSlugLimitDesc = `SELECT id, slug, title, message, forum_slug, user_nickname, created, votes
	FROM thread
	WHERE forum_slug = $1
	ORDER BY created DESC
	LIMIT $2;`

	getThreadsByForumSlugLimitSince = `SELECT id, slug, title, message, forum_slug, user_nickname, created, votes
	FROM thread
	WHERE forum_slug = $1 AND created >= $3::TEXT::TIMESTAMPTZ
	ORDER BY created
	LIMIT $2;`

	getThreadsByForumSlugLimitSinceDesc = `SELECT id, slug, title, message, forum_slug, user_nickname, created, votes
	FROM thread
	WHERE forum_slug = $1 AND created <= $3::TEXT::TIMESTAMPTZ
	ORDER BY created DESC
 	LIMIT $2;`

	checkThreadByIdOrSlug = `SELECT id, slug
	FROM thread
	WHERE slug = $1 OR id::TEXT = $1`

	updateThreadVotes = `UPDATE thread
	SET votes = votes + $1
	WHERE id = $2
	RETURNING id, slug, title, message, forum_slug, user_nickname, created, votes`

	updateThread = `UPDATE thread
	SET title = COALESCE($1, title), 
			message = COALESCE($2, message)
	WHERE id = $3
	RETURNING id, slug, title, message, forum_slug, user_nickname, created, votes`
)

func NewThreadRepo(conn *pgx.ConnPool) *Thread {
	return &Thread{
		conn: conn,
	}
}

type Thread struct {
	conn *pgx.ConnPool
}

func (t *Thread) CreateThread(data *thread.Create) (*thread.Thread, error) {
	tx, err := t.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var userID, forumID uint64

	if err := tx.QueryRow(getUserIdAndNicknameByNickname, data.UserNickname).Scan(&userID, &data.UserNickname); err != nil {
		return nil, err
	}

	if err := tx.QueryRow(getForumIdAndSlugBySlug, data.ForumSlug).Scan(&forumID, &data.ForumSlug); err != nil {
		return nil, err
	}

	received := &thread.Thread{}

	if data.Slug != nil {
		if err := tx.QueryRow(getThreadBySlug, data.Slug).Scan(&received.ID, &received.Slug, &received.Title, &received.Message, &received.ForumSlug, &received.UserNickname, &received.Created, &received.Votes); err == nil {
			return received, errors.New("threadAlreadyExists")
		}
	}

	if err := tx.QueryRow(createThread, data.Slug, data.Title, data.Message, forumID, data.ForumSlug, userID, data.UserNickname, data.Created).Scan(&received.ID, &received.Slug, &received.Title, &received.Message, &received.ForumSlug, &received.UserNickname, &received.Created, &received.Votes); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(updateForumThreads, forumID); err != nil {
		return nil, err
	}

	tx.Commit()
	return received, nil
}

func (t *Thread) GetThreads(slug string, limit *int, since *string, orderDesc bool) (*thread.Threads, error) {
	if err := t.conn.QueryRow(getForumSlugBySlug, slug).Scan(&slug); err != nil {
		return nil, err
	}

	threads := make(thread.Threads, 0, 0)
	var err error
	var rows *pgx.Rows

	if since == nil {
		if orderDesc {
			rows, err = t.conn.Query(getThreadsByForumSlugLimitDesc, slug, limit)
		} else {
			rows, err = t.conn.Query(getThreadsByForumSlugLimit, slug, limit)
		}
	} else {
		if orderDesc {
			rows, err = t.conn.Query(getThreadsByForumSlugLimitSinceDesc, slug, limit, since)
		} else {
			rows, err = t.conn.Query(getThreadsByForumSlugLimitSince, slug, limit, since)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var row thread.Thread
		rows.Scan(&row.ID, &row.Slug, &row.Title, &row.Message, &row.ForumSlug, &row.UserNickname, &row.Created, &row.Votes)
		threads = append(threads, row)
	}

	return &threads, nil
}

func (t *Thread) GetThread(slugOrId string) (*thread.Thread, error) {
	var received thread.Thread
	if err := t.conn.QueryRow(getThreadByIdOrSlug, slugOrId).
		Scan(&received.ID, &received.Slug, &received.Title, &received.Message, &received.ForumSlug, &received.UserNickname, &received.Created, &received.Votes); err != nil {
		return nil, err
	}

	return &received, nil
}

func (t *Thread) UpdateThread(data *thread.Update, slugOrId string) (*thread.Thread, error) {
	tx, err := t.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var threadID uint64
	var threadSlug string
	if err := tx.QueryRow(checkThreadByIdOrSlug, slugOrId).
		Scan(&threadID, &threadSlug); err != nil {
		return nil, err
	}

	var updated thread.Thread
	if err := tx.QueryRow(updateThread, data.Title, data.Message, threadID).
		Scan(&updated.ID, &updated.Slug, &updated.Title, &updated.Message, &updated.ForumSlug, &updated.UserNickname, &updated.Created, &updated.Votes); err != nil {
		return nil, err
	}

	tx.Commit()
	return &updated, nil
}
