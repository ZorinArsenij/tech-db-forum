package user

//go:generate easyjson user.go

//easyjson:json
type User struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
}

//easyjson:json
type Update struct {
	Email    *string `json:"email"`
	Fullname *string `json:"fullname"`
	About    *string `json:"about"`
}

type Info struct {
	ID       uint64
	Nickname string
}

//easyjson:json
type Users []User
