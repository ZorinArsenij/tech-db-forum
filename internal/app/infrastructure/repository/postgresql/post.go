package postgresql

import (
	"context"
	"errors"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
	"log"
	"strings"
	"time"

	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/forum"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/post"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/thread"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/user"

	"github.com/jackc/pgx"
)

const (
	getPostByIdAndThread             = "getPostByIdAndThread"
	getPostById                      = "getPostById"
	createPost                       = "createPost"
	createPostRoot                   = "createPostRoot"
	updatePost                       = "updatePost"
	getPostsFlat                     = "getPostsFlat"
	getPostsFlatLimit                = "getPostsFlatLimit"
	getPostsFlatLimitDesc            = "getPostsFlatLimitDesc"
	getPostsFlatLimitSince           = "getPostsFlatLimitSince"
	getPostsFlatLimitSinceDesc       = "getPostsFlatLimitSinceDesc"
	getPostsTree                     = "getPostsTree"
	getPostsTreeLimit                = "getPostsTreeLimit"
	getPostsTreeLimitDesc            = "getPostsTreeLimitDesc"
	getPostsTreeLimitSince           = "getPostsTreeLimitSince"
	getPostsTreeLimitSinceDesc       = "getPostsTreeLimitSinceDesc"
	getPostPath                      = "getPostPath"
	getPostRoot                      = "getPostRoot"
	getPostsParentTreeLimit          = "getPostsParentTreeLimit"
	getPostsParentTreeLimitDesc      = "getPostsParentTreeLimitDesc"
	getPostsParentTreeLimitSince     = "getPostsParentTreeLimitSince"
	getPostsParentTreeLimitSinceDesc = "getPostsParentTreeLimitSinceDesc"
	getPostsLimit                    = "getPostsLimit"
	getPostsLimitDesc                = "getPostsLimitDesc"
	getPostsLimitSince               = "getPostsLimitSince"
	getPostsLimitSinceDesc           = "getPostsLimitSinceDesc"
)

var postQueries = map[string]string{
	getPostByIdAndThread: `SELECT parents, root
	FROM post
	WHERE id = $1 AND thread_id = $2`,

	getPostById: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE id = $1;`,

	createPost: `INSERT INTO post (message, created, user_id, user_nickname, thread_id, forum_slug, parent, parents, root)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent`,

	createPostRoot: `INSERT INTO post (message, created, user_id, user_nickname, thread_id, forum_slug, parent, parents, root)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, (SELECT CURRVAL('post_id_seq')))
	RETURNING id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent`,

	updatePost: `UPDATE post
	SET message = COALESCE($1, message),
	is_edited = TRUE
	WHERE id = $2;`,

	getPostsFlat: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1`,

	getPostsFlatLimit: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1
	ORDER BY id, created
	LIMIT $2;`,

	getPostsFlatLimitDesc: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1
	ORDER BY id DESC, created
	LIMIT $2;`,

	getPostsFlatLimitSince: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1  AND id > $3
	ORDER BY id, created
	LIMIT $2;`,

	getPostsFlatLimitSinceDesc: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1  AND id < $3
	ORDER BY id DESC, created
	LIMIT $2;`,

	getPostsTree: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1`,

	getPostsTreeLimit: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1
	ORDER BY array_append(parents, id)
	LIMIT $2;`,

	getPostsTreeLimitDesc: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1
	ORDER BY array_append(parents, id) DESC
	LIMIT $2;`,

	getPostsTreeLimitSince: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1 AND array_append(parents, id) > $3
	ORDER BY array_append(parents, id)
	LIMIT $2;`,

	getPostsTreeLimitSinceDesc: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1 AND array_append(parents, id) < $3
	ORDER BY array_append(parents, id) DESC
	LIMIT $2;`,

	getPostPath: `SELECT array_append(parents, id) AS path
	FROM post
	WHERE id = $1;`,

	getPostRoot: `SELECT root::TEXT
	FROM post
	WHERE id = $1;`,

	getPostsParentTreeLimit: `SELECT p.id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post AS p
  	JOIN (
      SELECT id
      FROM post
      WHERE parent = 0
        AND thread_id = $1
      ORDER BY id
      LIMIT $2
	) AS s ON (p.root = s.id)
	ORDER BY root, array_append(p.parents, p.id)`,

	getPostsParentTreeLimitDesc: `SELECT p.id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post AS p
  	JOIN (
      SELECT id
      FROM post
      WHERE parent = 0
        AND thread_id = $1
      ORDER BY id DESC
      LIMIT $2
	) AS s ON (p.root = s.id)
	ORDER BY root DESC, array_append(p.parents, p.id)`,

	getPostsParentTreeLimitSince: `SELECT p.id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post AS p
  	JOIN (
      SELECT id
      FROM post
      WHERE parent = 0
        AND thread_id = $1
				AND root > $3
      ORDER BY id
      LIMIT $2
	) AS s ON (p.root = s.id)
	ORDER BY root, array_append(p.parents, p.id)`,

	getPostsParentTreeLimitSinceDesc: `SELECT p.id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post AS p
  	JOIN (
      SELECT id
      FROM post
      WHERE parent = 0
        AND thread_id = $1
				AND root < $3
      ORDER BY id DESC
      LIMIT $2
	) AS s ON (p.root = s.id)
	ORDER BY root DESC, array_append(p.parents, p.id)`,

	getPostsLimit: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1
	ORDER BY id
	LIMIT $2`,

	getPostsLimitDesc: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1
	ORDER BY id DESC
	LIMIT $2`,

	getPostsLimitSince: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1 AND id > $3
	ORDER BY id
	LIMIT $2`,

	getPostsLimitSinceDesc: `SELECT id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent
	FROM post
	WHERE thread_id = $1 AND id < $3
	ORDER BY id DESC
	LIMIT $2`,
}

func NewPostRepo(conn *pgx.ConnPool) *Post {
	return &Post{
		conn: conn,
	}
}

type Post struct {
	conn *pgx.ConnPool
}

func (p *Post) GetPost(id string, related map[string]bool) (*post.Info, error) {
	var info post.Info

	var post post.Post
	if err := p.conn.QueryRow(getPostById, id).Scan(&post.ID, &post.Message, &post.Created, &post.IsEdited, &post.UserNickname, &post.ThreadID, &post.ForumSlug, &post.Parent); err != nil {
		return nil, err
	}
	info.Post = post

	if value, exists := related["user"]; value && exists {
		var author user.User
		if err := p.conn.QueryRow(getUserByNickname, post.UserNickname).
			Scan(&author.Email, &author.Nickname, &author.Fullname, &author.About); err != nil {
			return nil, err
		}
		info.Author = &author
	}

	if value, exists := related["forum"]; value && exists {
		var relatedForum forum.Forum
		if err := p.conn.QueryRow(getForumBySlug, post.ForumSlug).
			Scan(&relatedForum.Slug, &relatedForum.Title, &relatedForum.Posts, &relatedForum.Threads, &relatedForum.UserNickname); err != nil {
			return nil, err
		}
		info.Forum = &relatedForum
	}

	if value, exists := related["thread"]; value && exists {
		var relatedThread thread.Thread
		if err := p.conn.QueryRow(getThreadById, post.ThreadID).
			Scan(&relatedThread.ID, &relatedThread.Slug, &relatedThread.Title, &relatedThread.Message, &relatedThread.ForumSlug, &relatedThread.UserNickname, &relatedThread.Created, &relatedThread.Votes); err != nil {
			return nil, err
		}
		info.Thread = &relatedThread
	}

	return &info, nil
}

func (p *Post) UpdatePost(data *post.Update) (*post.Post, error) {
	tx, err := p.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var received post.Post
	if err := tx.QueryRow(getPostById, data.ID).
		Scan(&received.ID, &received.Message, &received.Created, &received.IsEdited, &received.UserNickname, &received.ThreadID, &received.ForumSlug, &received.Parent); err != nil {
		return nil, err
	}

	if data.Message == nil || *data.Message == received.Message {
		return &received, nil
	}

	if _, err := tx.Exec(updatePost, data.Message, data.ID); err != nil {
		return nil, err
	}

	received.Message = *data.Message
	received.IsEdited = true
	tx.Commit()
	return &received, nil
}

func (p *Post) CreatePosts(data *post.PostsCreate, slugOrId string) (*post.Posts, error) {
	tx, err := p.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	batch := tx.BeginBatch()

	var threadID, userID uint64
	var forumSlug string

	createTime := time.Now()

	if err := tx.QueryRow(getThreadShortBySlugOrId, slugOrId).Scan(&threadID, &forumSlug); err != nil {
		batch.Close()
		return nil, err
	}

	postParentsCheckList := make([]int, 0, len(*data))
	userNicknamesCheckSet := treeset.NewWith(utils.StringComparator)
	for i, newPost := range *data {
		userNicknamesCheckSet.Add(newPost.UserNickname)
		if newPost.Parent != 0 {
			postParentsCheckList = append(postParentsCheckList, i)
			batch.Queue(getPostByIdAndThread, []interface{}{newPost.Parent, threadID}, nil, nil)
		}
	}

	userNicknamesList := userNicknamesCheckSet.Values()
	for _, nickname := range userNicknamesList {
		batch.Queue(getUserInfoByNickname, []interface{}{nickname}, nil, nil)
	}

	if err := batch.Send(context.Background(), nil); err != nil {
		batch.Close()
		log.Fatal(err)
	}

	for _, index := range postParentsCheckList {
		if err := batch.QueryRowResults().Scan(&(*data)[index].Parents, &(*data)[index].Root); err != nil {
			batch.Close()
			return nil, errors.New("postParentDoesNotExist")
		}
		(*data)[index].Parents = append((*data)[index].Parents, (*data)[index].Parent)
	}

	users := make(map[string]user.Info, len(userNicknamesList))
	for _, nickname := range userNicknamesList {
		info := user.Info{}
		if err := batch.QueryRowResults().Scan(&info.ID, &info.Email, &info.Nickname, &info.Fullname, &info.About); err != nil {
			batch.Close()
			return nil, err
		}
		users[strings.ToLower(nickname.(string))] = info
	}

	batch = tx.BeginBatch()

	for _, newPost := range *data {
		userID = users[strings.ToLower(newPost.UserNickname)].ID
		if newPost.Parent == 0 {
			batch.Queue(createPostRoot,
				[]interface{}{newPost.Message, createTime, userID, newPost.UserNickname, threadID, forumSlug, newPost.Parent, []int32{}},
				nil, nil)
		} else {
			batch.Queue(createPost,
				[]interface{}{newPost.Message, createTime, userID, newPost.UserNickname, threadID, forumSlug, newPost.Parent, newPost.Parents, newPost.Root},
				nil, nil)
		}
	}

	if err := batch.Send(context.Background(), nil); err != nil {
		batch.Close()
		log.Fatal(err)
	}

	posts := make(post.Posts, 0, len(*data))
	for range *data {
		var created post.Post
		if err := batch.QueryRowResults().
			Scan(&created.ID, &created.Message, &created.Created, &created.IsEdited, &created.UserNickname, &created.ThreadID, &created.ForumSlug, &created.Parent); err != nil {
			batch.Close()
			return nil, err
		}
		posts = append(posts, created)
	}

	batch = tx.BeginBatch()
	defer batch.Close()

	for _, info := range users {
		batch.Queue(createForumUser,
			[]interface{}{forumSlug, info.Email, info.Nickname, info.Fullname, info.About}, nil, nil)
	}

	if err := batch.Send(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	for range users {
		if _, err := batch.ExecResults(); err != nil {
			return nil, err
		}
	}

	if _, err := tx.Exec(updateForumPosts, len(*data), forumSlug); err != nil {
		return nil, err
	}

	tx.Commit()
	return &posts, nil
}

//func (p *Post) CreatePosts(data *post.PostsCreate, slugOrId string) (*post.Posts, error) {
//	if len(*data) == 0 {
//		return nil, nil
//	}
//
//	tx, err := p.conn.Begin()
//	if err != nil {
//		return nil, err
//	}
//	defer tx.Rollback()
//
//	var threadID, userID uint64
//	var forumSlug string
//
//	createTime := time.Now()
//	posts := make(post.Posts, 0, 0)
//
//	if err := tx.QueryRow(getThreadShortBySlugOrId, slugOrId).Scan(&threadID, &forumSlug); err != nil {
//		return nil, err
//	}
//
//	for _, newPost := range *data {
//		if err := tx.QueryRow(getUserIdAndNicknameByNickname, newPost.UserNickname).Scan(&userID, &newPost.UserNickname); err != nil {
//			return nil, err
//		}
//
//		parents := make([]int32, 0, 0)
//		var root int
//
//		if newPost.Parent != 0 {
//			parents = make([]int32, 0, 0)
//			if err := tx.QueryRow(getPostByIdAndThread, newPost.Parent, threadID).Scan(&parents, &root); err != nil {
//				return nil, errors.New("postParentDoesNotExist")
//			}
//
//			parents = append(parents, newPost.Parent)
//		}
//
//		var created post.Post
//		if newPost.Parent == 0 {
//			if err := tx.QueryRow(createPostRoot, newPost.Message, createTime, userID, newPost.UserNickname, threadID, forumSlug, newPost.Parent, parents).Scan(&created.ID, &created.Message, &created.Created, &created.IsEdited, &created.UserNickname, &created.ThreadID, &created.ForumSlug, &created.Parent); err != nil {
//				return nil, err
//			}
//		} else {
//			if err := tx.QueryRow(createPost, newPost.Message, createTime, userID, newPost.UserNickname, threadID, forumSlug, newPost.Parent, parents, root).Scan(&created.ID, &created.Message, &created.Created, &created.IsEdited, &created.UserNickname, &created.ThreadID, &created.ForumSlug, &created.Parent); err != nil {
//				return nil, err
//			}
//		}
//
//		posts = append(posts, created)
//	}
//
//	if _, err := tx.Exec(updateForumPosts, len(*data), forumSlug); err != nil {
//		return nil, err
//	}
//
//	tx.Commit()
//	return &posts, nil
//}

func (p *Post) GetPostsFlat(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error) {
	var threadID uint64
	if err := p.conn.QueryRow(getThreadShortBySlugOrId, slugOrId).
		Scan(&threadID, &slugOrId); err != nil {
		return nil, err
	}

	posts := make(post.Posts, 0)
	var err error
	var rows *pgx.Rows

	if since == nil {
		if orderDesc {
			rows, err = p.conn.Query(getPostsFlatLimitDesc, threadID, limit)
		} else {
			rows, err = p.conn.Query(getPostsFlatLimit, threadID, limit)
		}
	} else {
		if orderDesc {
			rows, err = p.conn.Query(getPostsFlatLimitSinceDesc, threadID, limit, since)
		} else {
			rows, err = p.conn.Query(getPostsFlatLimitSince, threadID, limit, since)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var row post.Post
		rows.Scan(&row.ID, &row.Message, &row.Created, &row.IsEdited, &row.UserNickname, &row.ThreadID, &row.ForumSlug, &row.Parent)
		posts = append(posts, row)
	}

	return &posts, nil
}

func (p *Post) GetPostsTree(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error) {
	var threadID uint64
	if err := p.conn.QueryRow(getThreadShortBySlugOrId, slugOrId).
		Scan(&threadID, &slugOrId); err != nil {
		return nil, err
	}

	posts := make(post.Posts, 0)
	var err error
	var rows *pgx.Rows

	if since != nil {
		parents := make([]int32, 0)
		_ = p.conn.QueryRow(getPostPath, since).Scan(&parents)
		if orderDesc {
			rows, err = p.conn.Query(getPostsTreeLimitSinceDesc, threadID, limit, parents)
		} else {
			rows, err = p.conn.Query(getPostsTreeLimitSince, threadID, limit, parents)
		}
	} else {
		if orderDesc {
			rows, err = p.conn.Query(getPostsTreeLimitDesc, threadID, limit)
		} else {
			rows, err = p.conn.Query(getPostsTreeLimit, threadID, limit)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var row post.Post
		rows.Scan(&row.ID, &row.Message, &row.Created, &row.IsEdited, &row.UserNickname, &row.ThreadID, &row.ForumSlug, &row.Parent)
		posts = append(posts, row)
	}

	return &posts, nil
}

func (p *Post) GetPostsParentTree(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error) {
	var threadID uint64
	if err := p.conn.QueryRow(getThreadShortBySlugOrId, slugOrId).
		Scan(&threadID, &slugOrId); err != nil {
		return nil, err
	}

	posts := make(post.Posts, 0)
	var err error
	var rows *pgx.Rows

	if since == nil {
		if orderDesc {
			rows, err = p.conn.Query(getPostsParentTreeLimitDesc, threadID, limit)
		} else {
			rows, err = p.conn.Query(getPostsParentTreeLimit, threadID, limit)
		}
	} else {
		_ = p.conn.QueryRow(getPostRoot, since).Scan(&since)
		if orderDesc {
			rows, err = p.conn.Query(getPostsParentTreeLimitSinceDesc, threadID, limit, since)
		} else {
			rows, err = p.conn.Query(getPostsParentTreeLimitSince, threadID, limit, since)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var row post.Post
		rows.Scan(&row.ID, &row.Message, &row.Created, &row.IsEdited, &row.UserNickname, &row.ThreadID, &row.ForumSlug, &row.Parent)
		posts = append(posts, row)
	}

	return &posts, nil
}

func (p *Post) GetPosts(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error) {
	var threadID uint64
	if err := p.conn.QueryRow(getThreadShortBySlugOrId, slugOrId).
		Scan(&threadID, &slugOrId); err != nil {
		return nil, err
	}

	posts := make(post.Posts, 0)
	var err error
	var rows *pgx.Rows

	if since == nil {
		if orderDesc {
			rows, err = p.conn.Query(getPostsLimitDesc, threadID, limit)
		} else {
			rows, err = p.conn.Query(getPostsLimit, threadID, limit)
		}
	} else {
		if orderDesc {
			rows, err = p.conn.Query(getPostsLimitSinceDesc, threadID, limit, since)
		} else {
			rows, err = p.conn.Query(getPostsLimitSince, threadID, limit, since)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var row post.Post
		rows.Scan(&row.ID, &row.Message, &row.Created, &row.IsEdited, &row.UserNickname, &row.ThreadID, &row.ForumSlug, &row.Parent)
		posts = append(posts, row)
	}

	return &posts, nil
}
