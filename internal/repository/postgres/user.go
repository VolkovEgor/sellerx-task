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

func (r *UserPg) ExistenceCheck(userId int) error {
	var tmp int
	query := fmt.Sprintf(`SELECT id FROM %s WHERE id = $1`, usersTable)
	row := r.db.QueryRow(query, userId)
	return row.Scan(&tmp)
}
