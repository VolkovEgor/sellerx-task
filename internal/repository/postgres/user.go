package postgres

import (
	"fmt"

	"github.com/VolkovEgor/sellerx-task/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserPg struct {
	db *sqlx.DB
}

func NewUserPg(db *sqlx.DB) *UserPg {
	return &UserPg{db: db}
}

func (r *UserPg) Create(user *model.User) (string, error) {
	var userId string
	query := fmt.Sprintf(
		`INSERT INTO %s (username, created_at)
		VALUES ($1, $2) RETURNING id`, usersTable)

	row := r.db.QueryRow(query, user.Username, user.CreatedAt)
	if err := row.Scan(&userId); err != nil {
		return "", err
	}

	return userId, nil
}

func (r *UserPg) ExistenceCheck(userId string) error {
	var tmp string
	query := fmt.Sprintf(`SELECT id FROM %s WHERE id = $1`, usersTable)
	row := r.db.QueryRow(query, userId)
	return row.Scan(&tmp)
}
