package thread

import "time"

//easyjson:json
type Thread struct {
	ID    uint64 `json:"id"`
	Votes int    `json:"votes"`
	Create
}

//easyjson:json
type Threads []Thread

//easyjson:json
type Create struct {
	Title        string     `json:"title"`
	Slug         *string    `json:"slug,omitempty"`
	Message      string     `json:"message"`
	Created      *time.Time `json:"created,omitempty"`
	UserNickname string     `json:"author"`
	ForumSlug    string     `json:"forum"`
}

//easyjson:json
type Update struct {
	Message *string `json:"message"`
	Title   *string `json:"title"`
}
