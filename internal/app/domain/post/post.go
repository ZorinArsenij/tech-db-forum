package post

//go:generate easyjson post.go

import (
	"time"

	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/forum"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/thread"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/user"
)

//easyjson:json
type Post struct {
	ID uint64 `json:"id"`

	Message  string    `json:"message"`
	Created  time.Time `json:"created"`
	IsEdited bool      `json:"isEdited"`

	UserNickname string `json:"author"`

	ThreadID uint64 `json:"thread"`

	ForumSlug string `json:"forum"`

	Parent int32 `json:"parent"`
}

//easyjson:json
type Posts []Post

//easyjson:json
type Create struct {
	Message      string `json:"message"`
	UserNickname string `json:"author"`
	Parent       int32  `json:"parent"`
}

//easyjson:json
type PostsCreate []Create

//easyjson:json
type Info struct {
	Author *user.User     `json:"author"`
	Forum  *forum.Forum   `json:"forum"`
	Post   Post           `json:"post"`
	Thread *thread.Thread `json:"thread"`
}

//easyjson:json
type Update struct {
	ID      string  `json:"-"`
	Message *string `json:"message"`
}

//type PostsQuery struct {
//	SlugOrId string
//	Limit *int
//	Since *string
//	OrderDesc bool
