package model

type Chat struct {
	Id              string   `json:"id" db:"id"`
	Name            string   `json:"name" db:"name"`
	Users           []string `json:"users" db:"users"`
	CreatedAt       int64    `json:"created_at" db:"created_at"`
	LastMessageTime int64    `json:"last_message_time" db:"last_message_time"`
}
