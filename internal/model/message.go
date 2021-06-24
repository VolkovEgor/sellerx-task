package model

type Message struct {
	Id        int    `json:"id" db:"id"`
	ChatId    int    `json:"chat_id" db:"chat_id"`
	AuthorId  int    `json:"author_id" db:"author_id"`
	Text      string `json:"text" db:"text"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
}
