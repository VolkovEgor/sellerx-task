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

func (r *MessagePg) Create(message *model.Message) (string, error) {
	var messageId string
	query := fmt.Sprintf(
		`INSERT INTO %s (chat_id, author_id, text, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`, messagesTable)

	row := r.db.QueryRow(query, message.ChatId, message.AuthorId, message.Text, message.CreatedAt)
	if err := row.Scan(&messageId); err != nil {
		return "", err
	}

	return messageId, nil
}

func (r *MessagePg) GetAllForChat(chatId string) ([]*model.Message, error) {
	messages := []*model.Message{}
	query := fmt.Sprintf(
		`SELECT id, chat_id, author_id, text, created_at
		FROM %s WHERE chat_id = $1
		ORDER BY created_at`, messagesTable)

	err := r.db.Select(&messages, query, chatId)
	return messages, err
}
