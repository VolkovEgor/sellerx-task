package postgres

import (
	"fmt"

	"github.com/VolkovEgor/sellerx-task/internal/model"

	"github.com/jmoiron/sqlx"
)

type ChatPg struct {
	db *sqlx.DB
}

func NewChatPg(db *sqlx.DB) *ChatPg {
	return &ChatPg{db: db}
}

func (r *ChatPg) Create(chat *model.Chat) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var chatId int
	query := fmt.Sprintf(
		`INSERT INTO %s (name, created_at)
		VALUES ($1, $2) RETURNING id`, chatsTable)

	row := r.db.QueryRow(query, chat.Name, chat.CreatedAt)
	if err := row.Scan(&chatId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	for _, userId := range chat.Users {
		var chatUserId int
		query = fmt.Sprintf(
			`INSERT INTO %s (chat_id, user_id)
			VALUES ($1, $2) RETURNING id`, chatUsersTable)

		row = r.db.QueryRow(query, chatId, userId)
		if err = row.Scan(&chatUserId); err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	return chatId, tx.Commit()
}

func (r *ChatPg) ExistenceCheck(chatId int) error {
	var tmp int
	query := fmt.Sprintf(`SELECT id FROM %s WHERE id = $1`, chatsTable)
	row := r.db.QueryRow(query, chatId)
	return row.Scan(&tmp)
}
