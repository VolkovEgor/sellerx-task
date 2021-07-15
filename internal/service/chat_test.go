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

func TestChatService_Create(t *testing.T) {
	type args struct {
		chat *model.Chat
	}

	type userMockBehavior func(r *mock_repositories.MockUser, userId string)
	type mockBehavior func(r *mock_repositories.MockChat, args args)

	tests := []struct {
		name     string
		input    args
		userMock []userMockBehavior
		mock     mockBehavior
		wantErr  bool
		want     string
	}{
		{
			name: "Ok",
			input: args{
				chat: &model.Chat{
					Name: "Chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
				},
			},
			userMock: []userMockBehavior{
				func(r *mock_repositories.MockUser, userId string) {
					r.EXPECT().ExistenceCheck(userId).Return(nil)
				},
				func(r *mock_repositories.MockUser, userId string) {
					r.EXPECT().ExistenceCheck(userId).Return(nil)
				},
			},
			mock: func(r *mock_repositories.MockChat, args args) {
				r.EXPECT().Create(args.chat).Return("00000000-0000-0000-0000-000000000001", nil)
			},
			want: "00000000-0000-0000-0000-000000000001",
		},
		{
			name: "Wrong chat name",
			input: args{
				chat: &model.Chat{
					Name: "",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
					},
				},
			},
			mock:     func(r *mock_repositories.MockChat, args args) {},
			userMock: []userMockBehavior{},
			wantErr:  true,
		},
		{
			name: "Empty user id",
			input: args{
				chat: &model.Chat{
					Name: "Chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"",
					},
				},
			},
			userMock: []userMockBehavior{
				func(r *mock_repositories.MockUser, userId string) {
					r.EXPECT().ExistenceCheck(userId).Return(nil)
				},
			},
			mock:    func(r *mock_repositories.MockChat, args args) {},
			wantErr: true,
		},
		{
			name: "No chat users",
			input: args{
				chat: &model.Chat{
					Name:  "Chat 1",
					Users: []string{},
				},
			},
			mock:     func(r *mock_repositories.MockChat, args args) {},
			userMock: []userMockBehavior{},
			wantErr:  true,
		},
		{
			name: "User not exists",
			input: args{
				chat: &model.Chat{
					Name: "Chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
				},
			},
			userMock: []userMockBehavior{
				func(r *mock_repositories.MockUser, userId string) {
					r.EXPECT().ExistenceCheck(userId).Return(nil)
				},
				func(r *mock_repositories.MockUser, userId string) {
					r.EXPECT().ExistenceCheck(userId).Return(sql.ErrNoRows)
				},
			},
			mock:    func(r *mock_repositories.MockChat, args args) {},
			wantErr: true,
		},
		{
			name: "Recurring users",
			input: args{
				chat: &model.Chat{
					Name: "Chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
						"00000000-0000-0000-0000-000000000001",
					},
				},
			},
			userMock: []userMockBehavior{
				func(r *mock_repositories.MockUser, userId string) {
					r.EXPECT().ExistenceCheck(userId).Return(nil)
				},
				func(r *mock_repositories.MockUser, userId string) {
					r.EXPECT().ExistenceCheck(userId).Return(nil)
				},
			},
			mock:    func(r *mock_repositories.MockChat, args args) {},
			wantErr: true,
		},
		{
			name: "Repo Error",
			input: args{
				chat: &model.Chat{
					Name: "Chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
					},
				},
			},
			userMock: []userMockBehavior{
				func(r *mock_repositories.MockUser, userId string) {
					r.EXPECT().ExistenceCheck(userId).Return(nil)
				},
			},
			mock: func(r *mock_repositories.MockChat, args args) {
				r.EXPECT().Create(args.chat).Return("", errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userRepo := mock_repositories.NewMockUser(c)
			repo := mock_repositories.NewMockChat(c)

			for i, userMock := range test.userMock {
				userMock(userRepo, test.input.chat.Users[i])
			}
			test.mock(repo, test.input)
			s := &ChatService{repo: repo, userRepo: userRepo}

			got, err := s.Create(test.input.chat)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestChatService_GetAllForUser(t *testing.T) {
	type args struct {
		user string
	}

	type userMockBehavior func(r *mock_repositories.MockUser, userId string)
	type mockBehavior func(r *mock_repositories.MockChat, args args)

	tests := []struct {
		name     string
		input    args
		userMock userMockBehavior
		mock     mockBehavior
		wantErr  bool
		want     []*model.Chat
	}{
		{
			name: "Ok",
			input: args{
				user: "00000000-0000-0000-0000-000000000001",
			},
			userMock: func(r *mock_repositories.MockUser, userId string) {
				r.EXPECT().ExistenceCheck(userId).Return(nil)
			},
			mock: func(r *mock_repositories.MockChat, args args) {
				data := []*model.Chat{
					{
						Id:   "00000000-0000-0000-0000-000000000001",
						Name: "Chat 1",
						Users: []string{
							"00000000-0000-0000-0000-000000000001",
							"00000000-0000-0000-0000-000000000002",
						},
						LastMessageTime: 200,
					},
					{
						Id:   "00000000-0000-0000-0000-000000000002",
						Name: "Chat 2",
						Users: []string{
							"00000000-0000-0000-0000-000000000003",
							"00000000-0000-0000-0000-000000000001",
						},
						LastMessageTime: 100,
					},
				}
				r.EXPECT().GetAllForUser(args.user).Return(data, nil)
			},
			want: []*model.Chat{
				{
					Id:   "00000000-0000-0000-0000-000000000001",
					Name: "Chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
					LastMessageTime: 200,
				},
				{
					Id:   "00000000-0000-0000-0000-000000000002",
					Name: "Chat 2",
					Users: []string{
						"00000000-0000-0000-0000-000000000003",
						"00000000-0000-0000-0000-000000000001",
					},
					LastMessageTime: 100,
				},
			},
		},
		{
			name: "Empty user id",
			input: args{
				user: "",
			},
			userMock: func(r *mock_repositories.MockUser, userId string) {},
			mock:     func(r *mock_repositories.MockChat, args args) {},
			wantErr:  true,
		},
		{
			name: "User not exists",
			input: args{
				user: "00000000-0000-0000-0000-000000000001",
			},
			userMock: func(r *mock_repositories.MockUser, userId string) {
				r.EXPECT().ExistenceCheck(userId).Return(sql.ErrNoRows)
			},
			mock:    func(r *mock_repositories.MockChat, args args) {},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userRepo := mock_repositories.NewMockUser(c)
			repo := mock_repositories.NewMockChat(c)

			test.userMock(userRepo, test.input.user)
			test.mock(repo, test.input)
			s := &ChatService{repo: repo, userRepo: userRepo}

			got, err := s.GetAllForUser(test.input.user)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}
