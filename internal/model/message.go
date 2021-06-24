package model

type Message struct {
	Id        string `json:"id" db:"id"`
	ChatId    string `json:"chat_id" db:"chat_id"`
	AuthorId  string `json:"author_id" db:"author_id"`
	Text      string `json:"text" db:"text"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
}
