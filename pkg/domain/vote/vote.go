package vote

//easyjson:json
type Vote struct {
	Rating       int    `json:"voice"`
	Voice        bool   `json:"-"`
	UserNickname string `json:"nickname"`
}
