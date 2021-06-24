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

func (r *ChatPg) GetAllForUser(userId int) ([]*model.Chat, error) {
	var chats []*model.Chat
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
	SELECT tmp.id, tmp.name, tmp.created_at, last_message_time
	FROM
	(	
		SELECT c.id, c.name, c.created_at, 
		-- Получение времени последнего сообщения в чате, если сообщений в чате нет, берется время создания чата
		(
			SELECT GREATEST(c.created_at, 
				(
					SELECT created_at FROM %s
					WHERE chat_id = c.id
					ORDER BY created_at DESC
					LIMIT 1
				)
			)
		) AS last_message_time
		FROM %s AS c INNER JOIN %s as cu ON c.id = cu.chat_id
		WHERE cu.user_id = $1
		ORDER BY last_message_time DESC 
	) AS tmp`, messagesTable, chatsTable, chatUsersTable)

	err = r.db.Select(&chats, query, userId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	for _, chat := range chats {
		query = fmt.Sprintf(`SELECT user_id FROM %s WHERE chat_id = $1`, chatUsersTable)
		err = r.db.Select(&chat.Users, query, chat.Id)

		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	return chats, tx.Commit()
}

func (r *ChatPg) ExistenceCheck(chatId int) error {
	var tmp int
	query := fmt.Sprintf(`SELECT id FROM %s WHERE id = $1`, chatsTable)
	row := r.db.QueryRow(query, chatId)
	return row.Scan(&tmp)
}
