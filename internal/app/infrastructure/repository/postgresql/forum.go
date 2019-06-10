package postgresql

import (
	"errors"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/forum"

	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/user"
	"github.com/jackc/pgx"
)

const (
	getForumBySlug              = "getForumBySlug"
	createForum                 = "createForum"
	getForumIdAndSlugBySlug     = "getForumIdAndSlugBySlug"
	updateForumThreads          = "updateForumThreads"
	updateForumPosts            = "updateForumPosts"
	getForumSlugBySlug          = "getForumSlugBySlug"
	getForumUsers               = "getForumUsers"
	getForumUsersLimit          = "getForumUsersLimit"
	getForumUsersLimitDesc      = "getForumUsersLimitDesc"
	getForumUsersLimitSince     = "getForumUsersLimitSince"
	getForumUsersLimitSinceDesc = "getForumUsersLimitSinceDesc"
)

var forumQueries = map[string]string{
	getForumBySlug: `SELECT slug, title, posts, threads, user_nickname
	FROM forum
	WHERE slug = $1;`,

	createForum: `INSERT INTO forum (slug, title, user_id, user_nickname) 
	VALUES ($1, $2, $3, $4)
	RETURNING slug, title, posts, threads, user_nickname;`,

	getForumIdAndSlugBySlug: `SELECT id, slug
	FROM forum
	WHERE slug = $1;`,

	updateForumThreads: `UPDATE forum
	SET threads = threads + 1
	WHERE id = $1;`,

	updateForumPosts: `UPDATE forum
	SET posts = posts + $1
	WHERE slug = $2;`,

	getForumSlugBySlug: `SELECT slug
	FROM forum
	WHERE slug = $1;`,

	getForumUsers: `SELECT c.nickname
	FROM (SELECT user_nickname AS nickname
					FROM thread
					WHERE forum_slug = $1
					GROUP BY user_nickname
					UNION
					SELECT user_nickname AS nickname
					FROM post
					WHERE forum_slug = $1
					GROUP BY user_nickname) AS u
	JOIN client AS c ON (c.nickname = u.nickname);`,

	getForumUsersLimit: `SELECT c.nickname
	FROM (SELECT user_nickname AS nickname
					FROM thread
					WHERE forum_slug = $1
					GROUP BY user_nickname
					UNION
					SELECT user_nickname AS nickname
					FROM post
					WHERE forum_slug = $1
					GROUP BY user_nickname) AS u
	JOIN client AS c ON (c.nickname = u.nickname)
	ORDER BY c.nickname
	LIMIT $2;`,

	getForumUsersLimitDesc: `SELECT c.nickname
	FROM (SELECT user_nickname AS nickname
					FROM thread
					WHERE forum_slug = $1
					GROUP BY user_nickname
					UNION
					SELECT user_nickname AS nickname
					FROM post
					WHERE forum_slug = $1
					GROUP BY user_nickname) AS u
	JOIN client AS c ON (c.nickname = u.nickname)
	ORDER BY c.nickname DESC
	LIMIT $2;`,

	getForumUsersLimitSince: `SELECT c.nickname
	FROM (SELECT user_nickname AS nickname
					FROM thread
					WHERE forum_slug = $1
					GROUP BY user_nickname
					UNION
					SELECT user_nickname AS nickname
					FROM post
					WHERE forum_slug = $1
					GROUP BY user_nickname) AS u
	JOIN client AS c ON (c.nickname = u.nickname)
	WHERE c.nickname > $3
	ORDER BY c.nickname
	LIMIT $2;`,

	getForumUsersLimitSinceDesc: `SELECT c.nickname
	FROM (SELECT user_nickname AS nickname
					FROM thread
					WHERE forum_slug = $1
					GROUP BY user_nickname
					UNION
					SELECT user_nickname AS nickname
					FROM post
					WHERE forum_slug = $1
					GROUP BY user_nickname) AS u
	JOIN client AS c ON (c.nickname = u.nickname)
	WHERE c.nickname < $3
	ORDER BY c.nickname DESC
	LIMIT $2;`,
}

func NewForumRepo(conn *pgx.ConnPool) *Forum {
	return &Forum{
		conn: conn,
	}
}

type Forum struct {
	conn *pgx.ConnPool
}

func (f *Forum) GetForum(slug string) (*forum.Forum, error) {
	received := &forum.Forum{}
	if err := f.conn.QueryRow(getForumBySlug, slug).Scan(&received.Slug, &received.Title, &received.Posts, &received.Threads, &received.UserNickname); err != nil {
		return nil, err
	}

	return received, nil
}

func (f *Forum) CreateForum(data *forum.Create) (*forum.Forum, error) {
	tx, err := f.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var id uint64

	if err := tx.QueryRow(getUserIdAndNicknameByNickname, data.UserNickname).Scan(&id, &data.UserNickname); err != nil {
		return nil, err
	}

	forum := &forum.Forum{}
	if err := tx.QueryRow(getForumBySlug, data.Slug).
		Scan(&forum.Slug, &forum.Title, &forum.Posts, &forum.Threads, &forum.UserNickname); err == nil {
		return forum, errors.New("ForumAlreadyExists")
	}

	if err := tx.QueryRow(createForum, data.Slug, data.Title, id, data.UserNickname).
		Scan(&forum.Slug, &forum.Title, &forum.Posts, &forum.Threads, &forum.UserNickname); err != nil {
		return nil, err
	}

	tx.Commit()
	return forum, nil
}

type UserQuery struct {
	Slug      string
	Limit     *int
	Since     *string
	OrderDesc bool
}

func (f *Forum) GetForumUsers(slug string, limit *int, since *string, orderDesc bool) (*user.Users, error) {
	if err := f.conn.QueryRow(getForumSlugBySlug, slug).Scan(&slug); err != nil {
		return nil, err
	}

	users := make(user.Users, 0)
	var err error
	var rows *pgx.Rows

	if since == nil {
		if orderDesc {
			rows, err = f.conn.Query(getForumUsersLimitDesc, slug, limit)
		} else {
			rows, err = f.conn.Query(getForumUsersLimit, slug, limit)
		}
	} else {
		if orderDesc {
			rows, err = f.conn.Query(getForumUsersLimitSinceDesc, slug, limit, since)
		} else {
			rows, err = f.conn.Query(getForumUsersLimitSince, slug, limit, since)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var received user.User
		var nickname string

		rows.Scan(&nickname)
		if err := f.conn.QueryRow(getUserByNickname, nickname).
			Scan(&received.Email, &received.Nickname, &received.Fullname, &received.About); err != nil {
			return nil, err
		}

		users = append(users, received)
	}

	return &users, nil
}
