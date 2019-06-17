package message

//go:generate easyjson message.go

//easyjson:json
type Message struct {
	Description string `json:"message"`
}
