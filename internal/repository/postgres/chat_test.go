package postgres

import (
	"errors"
	"fmt"
	"testing"

	"github.com/VolkovEgor/sellerx-task/internal/model"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestChatPg_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChatPg(db)

	type args struct {
		chat *model.Chat
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
				chat: &model.Chat{
					Name: "New test chat",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
					CreatedAt: 100,
				},
			},
			want: "00000000-0000-0000-0000-000000000001",
			mock: func(args args, id string) {
				mock.ExpectBegin()

				input := args.chat
				adRows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO chats").
					WithArgs(input.Name, input.CreatedAt).
					WillReturnRows(adRows)

				chatUsersId1 := "00000000-0000-0000-0000-000000000001"
				adRows = sqlmock.NewRows([]string{"id"}).AddRow(chatUsersId1)
				mock.ExpectQuery("INSERT INTO chat_users").
					WithArgs(id, input.Users[0]).
					WillReturnRows(adRows)

				chatUsersId2 := "00000000-0000-0000-0000-000000000002"
				adRows = sqlmock.NewRows([]string{"id"}).AddRow(chatUsersId2)
				mock.ExpectQuery("INSERT INTO chat_users").
					WithArgs(id, input.Users[1]).
					WillReturnRows(adRows)

				mock.ExpectCommit()
			},
		},
		{
			name: "Failed insert in chats table",
			input: args{
				chat: &model.Chat{
					Name: "New test chat",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
					CreatedAt: 100,
				},
			},
			wantErr: true,
			mock: func(args args, id string) {
				mock.ExpectBegin()
				input := args.chat
				mock.ExpectQuery("INSERT INTO chats").
					WithArgs(input.Name, input.CreatedAt).
					WillReturnError(fmt.Errorf("Some error"))
				mock.ExpectRollback()
			},
		},
		{
			name: "Failed insert in chat_users table",
			input: args{
				chat: &model.Chat{
					Name: "New test chat",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
					CreatedAt: 100,
				},
			},
			wantErr: true,
			mock: func(args args, id string) {
				mock.ExpectBegin()

				input := args.chat
				adRows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO chats").
					WithArgs(input.Name, input.CreatedAt).
					WillReturnRows(adRows)

				chatUsersId1 := "00000000-0000-0000-0000-000000000001"
				_ = sqlmock.NewRows([]string{"id"}).AddRow(chatUsersId1)
				mock.ExpectQuery("INSERT INTO chat_users").
					WithArgs(id, input.Users[0]).
					WillReturnError(fmt.Errorf("Some error"))

				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			got, err := r.Create(tt.input.chat)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestChatPg_GetAllForUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChatPg(db)

	type args struct {
		userId string
	}
	type mockBehavior func(args args)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    []*model.Chat
		wantErr bool
	}{
		{
			name: "Ok",
			input: args{
				userId: "00000000-0000-0000-0000-000000000001",
			},
			want: []*model.Chat{
				{
					Id:   "00000000-0000-0000-0000-000000000001",
					Name: "First chat",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
					CreatedAt:       100,
					LastMessageTime: 300,
				},
				{
					Id:   "00000000-0000-0000-0000-000000000002",
					Name: "Second chat",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000003",
					},
					CreatedAt:       200,
					LastMessageTime: 250,
				},
			},
			mock: func(args args) {
				chatId1 := "00000000-0000-0000-0000-000000000001"
				chatId2 := "00000000-0000-0000-0000-000000000002"
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "last_message_time"}).
					AddRow(chatId1, "First chat", 100, 300).
					AddRow(chatId2, "Second chat", 200, 250)
				mock.ExpectQuery("SELECT (.+) FROM (.+)").
					WithArgs(args.userId).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"created_at"}).
					AddRow("00000000-0000-0000-0000-000000000001").
					AddRow("00000000-0000-0000-0000-000000000002")
				mock.ExpectQuery("SELECT (.+) FROM chat_users").
					WithArgs(chatId1).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"created_at"}).
					AddRow("00000000-0000-0000-0000-000000000001").
					AddRow("00000000-0000-0000-0000-000000000003")
				mock.ExpectQuery("SELECT (.+) FROM chat_users").
					WithArgs(chatId2).WillReturnRows(rows)
			},
		},
		{
			name: "Failed insert in chats table",
			input: args{
				userId: "00000000-0000-0000-0000-000000000001",
			},
			wantErr: true,
			mock: func(args args) {
				mock.ExpectQuery("SELECT (.+) FROM (.+)").
					WithArgs(args.userId).
					WillReturnError(fmt.Errorf("Some error"))
			},
		},
		{
			name: "Failed insert in chat_users table",
			input: args{
				userId: "00000000-0000-0000-0000-000000000001",
			},
			wantErr: true,
			mock: func(args args) {
				chatId1 := "00000000-0000-0000-0000-000000000001"
				chatId2 := "00000000-0000-0000-0000-000000000002"
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "last_message_time"}).
					AddRow(chatId1, "First chat", 100, 300).
					AddRow(chatId2, "Second chat", 200, 250)
				mock.ExpectQuery("SELECT (.+) FROM (.+)").
					WithArgs(args.userId).WillReturnRows(rows)

				mock.ExpectQuery("SELECT (.+) FROM chat_users").
					WithArgs(chatId1).
					WillReturnError(fmt.Errorf("Some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.GetAllForUser(tt.input.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestChatPg_ExistenceCheck(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChatPg(db)

	type args struct {
		chatId string
	}
	type mockBehavior func(args args)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    error
		wantErr bool
	}{
		{
			name: "Ok",
			input: args{
				chatId: "00000000-0000-0000-0000-000000000001",
			},
			want: nil,
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("00000000-0000-0000-0000-000000000001")
				mock.ExpectQuery("SELECT (.+) FROM chats").
					WithArgs(args.chatId).WillReturnRows(rows)
			},
		},
		{
			name: "Failed check",
			input: args{
				chatId: "00000000-0000-0000-0000-000000000001",
			},
			want: errors.New("Some error"),
			mock: func(args args) {
				mock.ExpectQuery("SELECT (.+) FROM chats").
					WithArgs(args.chatId).
					WillReturnError(errors.New("Some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			err := r.ExistenceCheck(tt.input.chatId)
			assert.Equal(t, tt.want, err)
		})
	}
}
