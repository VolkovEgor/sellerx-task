package service

import (
	"errors"
	"testing"

	"github.com/VolkovEgor/sellerx-task/internal/model"
	mock_repositories "github.com/VolkovEgor/sellerx-task/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Create(t *testing.T) {
	type args struct {
		user *model.User
	}
	type mockBehavior func(r *mock_repositories.MockUser, args args)

	tests := []struct {
		name    string
		input   args
		mock    mockBehavior
		wantErr bool
		want    string
	}{
		{
			name: "Ok",
			input: args{
				user: &model.User{
					Username: "User 1",
				},
			},
			mock: func(r *mock_repositories.MockUser, args args) {
				r.EXPECT().Create(args.user).Return("00000000-0000-0000-0000-000000000001", nil)
			},
			want: "00000000-0000-0000-0000-000000000001",
		},
		{
			name: "Wrong username",
			input: args{
				user: &model.User{
					Username: "",
				},
			},
			mock:    func(r *mock_repositories.MockUser, args args) {},
			wantErr: true,
		},
		{
			name: "Repo Error",
			input: args{
				user: &model.User{
					Username: "User 1",
				},
			},
			mock: func(r *mock_repositories.MockUser, args args) {
				r.EXPECT().Create(args.user).Return("", errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repositories.NewMockUser(c)
			test.mock(repo, test.input)
			s := &UserService{repo: repo}

			got, err := s.Create(test.input.user)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}
