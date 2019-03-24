package forum

//easyjson:json
type Forum struct {
	ID           uint64 `json:"-"`
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	Threads      int    `json:"threads"`
	Posts        int64  `json:"posts"`
	UserID       uint64 `json:"-"`
	UserNickname string `json:"user"`
}

//easyjson:json
type Create struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	UserNickname string `json:"user"`
}
