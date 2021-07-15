package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/VolkovEgor/sellerx-task/internal/model"
	mock_repositories "github.com/VolkovEgor/sellerx-task/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMessageService_Create(t *testing.T) {
	type args struct {
		message *model.Message
	}

	type chatMockBehavior func(r *mock_repositories.MockChat, chatId string)
	type userMockBehavior func(r *mock_repositories.MockUser, userId string)
	type mockBehavior func(r *mock_repositories.MockMessage, args args)

	tests := []struct {
		name     string
		input    args
		chatMock chatMockBehavior
		userMock userMockBehavior
		mock     mockBehavior
		wantErr  bool
		want     string
	}{
		{
			name: "Ok",
			input: args{
				message: &model.Message{
					ChatId:   "00000000-0000-0000-0000-000000000001",
					AuthorId: "00000000-0000-0000-0000-000000000001",
					Text:     "Message 1",
				},
			},
			chatMock: func(r *mock_repositories.MockChat, chatId string) {
				chat := &model.Chat{
					Id:   "00000000-0000-0000-0000-000000000001",
					Name: "Chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
					CreatedAt:       100,
					LastMessageTime: 200,
				}
				r.EXPECT().GetById(chatId).Return(chat, nil)
			},
			userMock: func(r *mock_repositories.MockUser, userId string) {
				r.EXPECT().ExistenceCheck(userId).Return(nil)
			},
			mock: func(r *mock_repositories.MockMessage, args args) {
				r.EXPECT().Create(args.message).Return("00000000-0000-0000-0000-000000000001", nil)
			},
			want: "00000000-0000-0000-0000-000000000001",
		},
		{
			name: "Wrong message text",
			input: args{
				message: &model.Message{
					ChatId:   "00000000-0000-0000-0000-000000000001",
					AuthorId: "00000000-0000-0000-0000-000000000001",
					Text:     "",
				},
			},
			chatMock: func(r *mock_repositories.MockChat, chatId string) {},
			userMock: func(r *mock_repositories.MockUser, userId string) {},
			mock:     func(r *mock_repositories.MockMessage, args args) {},
			wantErr:  true,
		},
		{
			name: "Wrong chat id",
			input: args{
				message: &model.Message{
					ChatId:   "",
					AuthorId: "00000000-0000-0000-0000-000000000001",
					Text:     "Message 1",
				},
			},
			chatMock: func(r *mock_repositories.MockChat, chatId string) {},
			userMock: func(r *mock_repositories.MockUser, userId string) {},
			mock:     func(r *mock_repositories.MockMessage, args args) {},
			wantErr:  true,
		},
		{
			name: "Wrong author id",
			input: args{
				message: &model.Message{
					ChatId:   "00000000-0000-0000-0000-000000000001",
					AuthorId: "",
					Text:     "Message 1",
				},
			},
			chatMock: func(r *mock_repositories.MockChat, chatId string) {},
			userMock: func(r *mock_repositories.MockUser, userId string) {},
			mock:     func(r *mock_repositories.MockMessage, args args) {},
			wantErr:  true,
		},
		{
			name: "Chat not exists",
			input: args{
				message: &model.Message{
					ChatId:   "00000000-0000-0000-0000-000000000001",
					AuthorId: "00000000-0000-0000-0000-000000000001",
					Text:     "Message 1",
				},
			},
			chatMock: func(r *mock_repositories.MockChat, chatId string) {
				r.EXPECT().GetById(chatId).Return(nil, sql.ErrNoRows)
			},
			userMock: func(r *mock_repositories.MockUser, userId string) {},
			mock:     func(r *mock_repositories.MockMessage, args args) {},
			wantErr:  true,
		},
		{
			name: "User not exists",
			input: args{
				message: &model.Message{
					ChatId:   "00000000-0000-0000-0000-000000000001",
					AuthorId: "00000000-0000-0000-0000-000000000001",
					Text:     "Message 1",
				},
			},
			chatMock: func(r *mock_repositories.MockChat, chatId string) {
				chat := &model.Chat{
					Id:   "00000000-0000-0000-0000-000000000001",
					Name: "Chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
					CreatedAt:       100,
					LastMessageTime: 200,
				}
				r.EXPECT().GetById(chatId).Return(chat, nil)
			},
			userMock: func(r *mock_repositories.MockUser, userId string) {
				r.EXPECT().ExistenceCheck(userId).Return(sql.ErrNoRows)
			},
			mock:    func(r *mock_repositories.MockMessage, args args) {},
			wantErr: true,
		},
		{
			name: "Ok",
			input: args{
				message: &model.Message{
					ChatId:   "00000000-0000-0000-0000-000000000001",
					AuthorId: "00000000-0000-0000-0000-000000000001",
					Text:     "Message 1",
				},
			},
			chatMock: func(r *mock_repositories.MockChat, chatId string) {
				chat := &model.Chat{
					Id:   "00000000-0000-0000-0000-000000000001",
					Name: "Chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
					CreatedAt:       int64(^uint64(0) >> 1), // max int64 value
					LastMessageTime: 200,
				}
				r.EXPECT().GetById(chatId).Return(chat, nil)
			},
			userMock: func(r *mock_repositories.MockUser, userId string) {
				r.EXPECT().ExistenceCheck(userId).Return(nil)
			},
			mock:    func(r *mock_repositories.MockMessage, args args) {},
			wantErr: true,
		},

		{
			name: "Repo Error",
			input: args{
				message: &model.Message{
					ChatId:   "00000000-0000-0000-0000-000000000001",
					AuthorId: "00000000-0000-0000-0000-000000000001",
					Text:     "Message 1",
				},
			},
			chatMock: func(r *mock_repositories.MockChat, chatId string) {
				chat := &model.Chat{
					Id:   "00000000-0000-0000-0000-000000000001",
					Name: "Chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
					CreatedAt:       100,
					LastMessageTime: 200,
				}
				r.EXPECT().GetById(chatId).Return(chat, nil)
			},
			userMock: func(r *mock_repositories.MockUser, userId string) {
				r.EXPECT().ExistenceCheck(userId).Return(nil)
			},
			mock: func(r *mock_repositories.MockMessage, args args) {
				r.EXPECT().Create(args.message).Return("", errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userRepo := mock_repositories.NewMockUser(c)
			chatRepo := mock_repositories.NewMockChat(c)
			repo := mock_repositories.NewMockMessage(c)

			test.userMock(userRepo, test.input.message.AuthorId)
			test.chatMock(chatRepo, test.input.message.ChatId)
			test.mock(repo, test.input)
			s := &MessageService{
				repo:     repo,
				userRepo: userRepo,
				chatRepo: chatRepo,
			}

			got, err := s.Create(test.input.message)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestMessageService_GetAllForChat(t *testing.T) {
	type args struct {
		chat string
	}

	type chatMockBehavior func(r *mock_repositories.MockChat, chatId string)
	type mockBehavior func(r *mock_repositories.MockMessage, args args)

	tests := []struct {
		name     string
		input    args
		chatMock chatMockBehavior
		mock     mockBehavior
		wantErr  bool
		want     []*model.Message
	}{
		{
			name: "Ok",
			input: args{
				chat: "00000000-0000-0000-0000-000000000001",
			},
			chatMock: func(r *mock_repositories.MockChat, chatId string) {
				r.EXPECT().ExistenceCheck(chatId).Return(nil)
			},
			mock: func(r *mock_repositories.MockMessage, args args) {
				data := []*model.Message{
					{
						Id:        "00000000-0000-0000-0000-000000000001",
						ChatId:    "00000000-0000-0000-0000-000000000001",
						AuthorId:  "00000000-0000-0000-0000-000000000001",
						Text:      "Message 2",
						CreatedAt: 300,
					},
					{
						Id:        "00000000-0000-0000-0000-000000000001",
						ChatId:    "00000000-0000-0000-0000-000000000001",
						AuthorId:  "00000000-0000-0000-0000-000000000001",
						Text:      "Message 1",
						CreatedAt: 200,
					},
				}
				r.EXPECT().GetAllForChat(args.chat).Return(data, nil)
			},
			want: []*model.Message{
				{
					Id:        "00000000-0000-0000-0000-000000000001",
					ChatId:    "00000000-0000-0000-0000-000000000001",
					AuthorId:  "00000000-0000-0000-0000-000000000001",
					Text:      "Message 2",
					CreatedAt: 300,
				},
				{
					Id:        "00000000-0000-0000-0000-000000000001",
					ChatId:    "00000000-0000-0000-0000-000000000001",
					AuthorId:  "00000000-0000-0000-0000-000000000001",
					Text:      "Message 1",
					CreatedAt: 200,
				},
			},
		},
		{
			name: "Empty chat id",
			input: args{
				chat: "",
			},
			chatMock: func(r *mock_repositories.MockChat, userId string) {},
			mock:     func(r *mock_repositories.MockMessage, args args) {},
			wantErr:  true,
		},
		{
			name: "Chat not exists",
			input: args{
				chat: "00000000-0000-0000-0000-000000000001",
			},
			chatMock: func(r *mock_repositories.MockChat, chatId string) {
				r.EXPECT().ExistenceCheck(chatId).Return(sql.ErrNoRows)
			},
			mock:    func(r *mock_repositories.MockMessage, args args) {},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userRepo := mock_repositories.NewMockUser(c)
			chatRepo := mock_repositories.NewMockChat(c)
			repo := mock_repositories.NewMockMessage(c)

			test.chatMock(chatRepo, test.input.chat)
			test.mock(repo, test.input)
			s := &MessageService{repo: repo, userRepo: userRepo, chatRepo: chatRepo}

			got, err := s.GetAllForChat(test.input.chat)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}
