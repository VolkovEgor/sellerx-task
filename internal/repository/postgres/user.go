package postgres

import (
	"fmt"

	"github.com/VolkovEgor/sellerx-task/internal/model"

	"github.com/jmoiron/sqlx"
)

const (
	NumElementsInPage = 10
)

type UserPg struct {
	db *sqlx.DB
}

func NewUserPg(db *sqlx.DB) *UserPg {
	return &UserPg{db: db}
}

func (r *UserPg) Create(user *model.User) (int, error) {
	var userId int
	query := fmt.Sprintf(
		`INSERT INTO %s (username, created_at)
		VALUES ($1, $2) RETURNING id`, usersTable)

	row := r.db.QueryRow(query, user.Username, user.CreatedAt)
	if err := row.Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}
