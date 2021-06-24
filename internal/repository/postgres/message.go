package postgres

import (
	"fmt"

	"github.com/VolkovEgor/sellerx-task/internal/model"

	"github.com/jmoiron/sqlx"
)

type MessagePg struct {
	db *sqlx.DB
}

func NewMessagePg(db *sqlx.DB) *MessagePg {
	return &MessagePg{db: db}
}

func (r *MessagePg) Create(message *model.Message) (int, error) {
	var messageId int
	query := fmt.Sprintf(
		`INSERT INTO %s (chat_id, author_id, text, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`, messagesTable)

	row := r.db.QueryRow(query, message.ChatId, message.AuthorId, message.Text, message.CreatedAt)
	if err := row.Scan(&messageId); err != nil {
		return 0, err
	}

	return messageId, nil
}
