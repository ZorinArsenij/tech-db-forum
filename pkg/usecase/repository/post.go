package repository

import "github.com/ZorinArsenij/tech-db-forum/pkg/domain/post"

type Post interface {
	GetPost(id string, related map[string]bool) (*post.Info, error)
	UpdatePost(data *post.Update) (*post.Post, error)
	CreatePosts(data *post.PostsCreate, slugOrId string) (*post.Posts, error)
	GetPosts(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error)
	GetPostsParentTree(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error)
	GetPostsTree(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error)
	GetPostsFlat(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error)
}
