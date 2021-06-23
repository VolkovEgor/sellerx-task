package model

type User struct {
	Id        int    `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
}
