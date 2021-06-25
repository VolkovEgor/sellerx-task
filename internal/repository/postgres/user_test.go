package postgres

import (
	"fmt"
	"testing"

	"github.com/VolkovEgor/sellerx-task/internal/model"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUserPg_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserPg(db)

	type args struct {
		user *model.User
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
				user: &model.User{
					Username:  "New Test User",
					CreatedAt: 100,
				},
			},
			want: "00000000-0000-0000-0000-000000000001",
			mock: func(args args, id string) {
				input := args.user
				adRows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(input.Username, input.CreatedAt).
					WillReturnRows(adRows)
			},
		},
		{
			name: "Failed insert",
			input: args{
				user: &model.User{
					Username:  "New Test User",
					CreatedAt: 100,
				},
			},
			wantErr: true,
			mock: func(args args, id string) {
				input := args.user
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(input.Username, input.CreatedAt).
					WillReturnError(fmt.Errorf("Some error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			got, err := r.Create(tt.input.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
