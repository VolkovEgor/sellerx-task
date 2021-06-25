package postgres

import (
	"fmt"
	"testing"

	"github.com/VolkovEgor/sellerx-task/internal/model"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestMessagePg_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewMessagePg(db)

	type args struct {
		message *model.Message
	}
	type mockBehavior func(args args, id string)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    string
		wantErr bool
	}{
		{
			name: "Ok",
			input: args{
				message: &model.Message{
					ChatId:    "00000000-0000-0000-0000-000000000001",
					AuthorId:  "00000000-0000-0000-0000-000000000001",
					Text:      "Some text",
					CreatedAt: 100,
				},
			},
			want: "00000000-0000-0000-0000-000000000001",
			mock: func(args args, id string) {
				input := args.message
				adRows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO messages").
					WithArgs(input.ChatId, input.AuthorId, input.Text, input.CreatedAt).
					WillReturnRows(adRows)
			},
		},
		{
			name: "Failed insert",
			input: args{
				message: &model.Message{
					ChatId:    "00000000-0000-0000-0000-000000000001",
					AuthorId:  "00000000-0000-0000-0000-000000000001",
					Text:      "Some text",
					CreatedAt: 100,
				},
			},
			wantErr: true,
			mock: func(args args, id string) {
				input := args.message
				mock.ExpectQuery("INSERT INTO messages").
					WithArgs(input.ChatId, input.AuthorId, input.Text, input.CreatedAt).
					WillReturnError(fmt.Errorf("Some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			got, err := r.Create(tt.input.message)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMessagePg_GetAllForChat(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewMessagePg(db)

	type args struct {
		chatId string
	}
	type mockBehavior func(args args)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    []*model.Message
		wantErr bool
	}{
		{
			name: "Ok",
			input: args{
				chatId: "00000000-0000-0000-0000-000000000001",
			},
			want: []*model.Message{
				{
					Id:        "00000000-0000-0000-0000-000000000001",
					ChatId:    "00000000-0000-0000-0000-000000000001",
					AuthorId:  "00000000-0000-0000-0000-000000000001",
					Text:      "First message text",
					CreatedAt: 100,
				},
				{
					Id:        "00000000-0000-0000-0000-000000000002",
					ChatId:    "00000000-0000-0000-0000-000000000001",
					AuthorId:  "00000000-0000-0000-0000-000000000002",
					Text:      "Second message text",
					CreatedAt: 200,
				},
			},
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "chat_id", "author_id", "text", "created_at"}).
					AddRow(
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000001",
						"First message text", 100).
					AddRow(
						"00000000-0000-0000-0000-000000000002",
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
						"Second message text", 200)
				mock.ExpectQuery("SELECT (.+) FROM messages").
					WithArgs(args.chatId).WillReturnRows(rows)
			},
		},
		{
			name: "Failed insert",
			input: args{
				chatId: "00000000-0000-0000-0000-000000000001",
			},
			wantErr: true,
			mock: func(args args) {
				mock.ExpectQuery("SELECT (.+) FROM messages").
					WithArgs(args.chatId).
					WillReturnError(fmt.Errorf("Some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.GetAllForChat(tt.input.chatId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
