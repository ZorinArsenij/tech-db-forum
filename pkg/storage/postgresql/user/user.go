package user

//easyjson:json
type User struct {
	ID       uint64 `json:"id,omitempty"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
}

//easyjson:json
type UserUpdate struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
}

//easyjson:json
type Users []User
