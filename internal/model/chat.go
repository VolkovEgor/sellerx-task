package model

type Chat struct {
	Id        int      `json:"id" db:"id"`
	Name      string   `json:"name" db:"name"`
	Users     []string `json:"users"`
	CreatedAt int64    `json:"created_at" db:"created_at"`
}
