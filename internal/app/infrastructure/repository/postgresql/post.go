package postgresql

import (
	"context"
	"errors"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/infrastructure/repository/postgresql/cluster"
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

	createPost: `INSERT INTO post (message, created, user_nickname, thread_id, forum_slug, parent, parents, root)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent`,

	createPostRoot: `INSERT INTO post (message, created, user_nickname, thread_id, forum_slug, parent, parents, root)
	VALUES ($1, $2, $3, $4, $5, $6, $7, (SELECT CURRVAL('post_id_seq')))
	RETURNING id, message, created, is_edited, user_nickname, thread_id, forum_slug, parent`,

	updatePost: `UPDATE post
	SET message = COALESCE($1, message),
	is_edited = TRUE
	WHERE id = $2;`,

	// Index???
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

var (
	CurrentPostNumber = 0
)

const (
	ClusteringStep = 1500000
)

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
			Scan(&author.Nickname, &author.Email, &author.Fullname, &author.About); err != nil {
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

func getUsersBatch(tx *pgx.Tx, data *post.PostsCreate) (*map[string]user.Info, error) {
	batch := tx.BeginBatch()
	defer batch.Close()

	nicknamesSet := treeset.NewWith(utils.StringComparator)
	for _, newPost := range *data {
		nicknamesSet.Add(newPost.UserNickname)
	}

	nicknames := nicknamesSet.Values()
	for _, nickname := range nicknames {
		batch.Queue(getUserByNickname, []interface{}{nickname}, nil, nil)
	}

	if err := batch.Send(context.Background(), nil); err != nil {
		return nil, err
	}

	users := make(map[string]user.Info, len(nicknames))
	for _, nickname := range nicknames {
		info := user.Info{}
		if err := batch.QueryRowResults().Scan(&info.Nickname, &info.Email, &info.Fullname, &info.About); err != nil {
			return nil, err
		}
		users[strings.ToLower(nickname.(string))] = info
	}

	return &users, nil
}

func getPostParentsBatch(tx *pgx.Tx, data *post.PostsCreate, threadID uint64) error {
	batch := tx.BeginBatch()
	defer batch.Close()

	postParentsCheckList := make([]int, 0, len(*data))
	for i, newPost := range *data {
		if newPost.Parent != 0 {
			postParentsCheckList = append(postParentsCheckList, i)
			batch.Queue(getPostByIdAndThread, []interface{}{newPost.Parent, threadID}, nil, nil)
		}
	}

	if err := batch.Send(context.Background(), nil); err != nil {
		return err
	}

	for _, index := range postParentsCheckList {
		if err := batch.QueryRowResults().Scan(&(*data)[index].Parents, &(*data)[index].Root); err != nil {
			return errors.New("postParentDoesNotExist")
		}
		(*data)[index].Parents = append((*data)[index].Parents, (*data)[index].Parent)
	}

	return nil
}

func createPostsBatch(tx *pgx.Tx, data *post.PostsCreate, threadID uint64, forumSlug string) (*post.Posts, error) {
	batch := tx.BeginBatch()
	defer batch.Close()

	createTime := time.Now()

	for _, newPost := range *data {
		if newPost.Parent == 0 {
			batch.Queue(createPostRoot,
				[]interface{}{newPost.Message, createTime, newPost.UserNickname, threadID, forumSlug, newPost.Parent, []int32{}},
				nil, nil)
		} else {
			batch.Queue(createPost,
				[]interface{}{newPost.Message, createTime, newPost.UserNickname, threadID, forumSlug, newPost.Parent, newPost.Parents, newPost.Root},
				nil, nil)
		}
	}

	if err := batch.Send(context.Background(), nil); err != nil {
		return nil, err
	}

	posts := make(post.Posts, 0, len(*data))
	for range *data {
		var created post.Post
		if err := batch.QueryRowResults().
			Scan(&created.ID, &created.Message, &created.Created, &created.IsEdited, &created.UserNickname, &created.ThreadID, &created.ForumSlug, &created.Parent); err != nil {
			return nil, err
		}
		posts = append(posts, created)
	}

	return &posts, nil
}

func createForumUsers(conn *pgx.ConnPool, forumSlug string, users *map[string]user.Info) error {
	for _, info := range *users {
		if _, err := conn.Exec(createForumUser, forumSlug, info.Email, info.Nickname, info.Fullname, info.About); err != nil {
			return err
		}
	}

	return nil
}

func (p *Post) CreatePosts(data *post.PostsCreate, slugOrId string) (*post.Posts, error) {
	tx, err := p.conn.Begin()
	if err != nil {
		log.Println("[Failed] creating transaction. Error:", err)
		return nil, err
	}
	defer tx.Rollback()

	var threadID uint64
	var forumSlug string

	if err := tx.QueryRow(getThreadShortBySlugOrId, slugOrId).Scan(&threadID, &forumSlug); err != nil {
		log.Println("[Failed] get threadId by forum slug or id. Error:", err)
		return nil, err
	}

	users, err := getUsersBatch(tx, data)
	if err != nil {
		log.Println("[Failed] getting users using batch. Error:", err)
		return nil, err
	}

	if err := getPostParentsBatch(tx, data, threadID); err != nil {
		log.Println("[Failed] getting posts parents. Error:", err)
		return nil, err
	}

	posts, err := createPostsBatch(tx, data, threadID, forumSlug)
	if err != nil {
		log.Println("[Failed] creating posts. Error:", err)
		return nil, err
	}

	if _, err := tx.Exec(updateForumPosts, len(*data), forumSlug); err != nil {
		log.Println("[Failed] updating forum posts. Error:", err)
		return nil, err
	}
	tx.Commit()

	if err := createForumUsers(p.conn, forumSlug, users); err != nil {
		log.Println("[Failed] creating forum users. Error:", err)
		return nil, err
	}

	CurrentPostNumber += len(*data)
	if CurrentPostNumber >= ClusteringStep {
		if err := cluster.CreateClusters(p.conn, "build/schema/1_cluster.sql"); err != nil {
			log.Fatal("[Failed] creating clusters. Error:", err)
		}
	}

	return posts, nil
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
