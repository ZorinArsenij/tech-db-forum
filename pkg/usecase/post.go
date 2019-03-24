package usecase

import (
	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/post"
	"github.com/ZorinArsenij/tech-db-forum/pkg/usecase/repository"
)

func NewPostInteractor(repo repository.Post) *PostInteractor {
	return &PostInteractor{
		repository: repo,
	}
}

type PostInteractor struct {
	repository repository.Post
}

func (i *PostInteractor) GetPost(id string, related map[string]bool) (*post.Info, error) {
	return i.repository.GetPost(id, related)
}

func (i *PostInteractor) UpdatePost(data *post.Update) (*post.Post, error) {
	return i.repository.UpdatePost(data)
}

func (i *PostInteractor) CreatePosts(data *post.PostsCreate, slugOrId string) (*post.Posts, error) {
	return i.repository.CreatePosts(data, slugOrId)
}

func (i *PostInteractor) GetPosts(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error) {
	return i.repository.GetPosts(slugOrId, limit, since, orderDesc)
}

func (i *PostInteractor) GetPostsParentTree(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error) {
	return i.repository.GetPostsParentTree(slugOrId, limit, since, orderDesc)
}

func (i *PostInteractor) GetPostsTree(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error) {
	return i.repository.GetPostsTree(slugOrId, limit, since, orderDesc)
}

func (i *PostInteractor) GetPostsFlat(slugOrId string, limit *int, since *string, orderDesc bool) (*post.Posts, error) {
	return i.repository.GetPostsFlat(slugOrId, limit, since, orderDesc)
}
