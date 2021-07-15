package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	errMes "github.com/VolkovEgor/sellerx-task/internal/error_message"
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/service"
	mock_service "github.com/VolkovEgor/sellerx-task/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func getBodyForError(code int, err string) string {
	return `{"error_code":` + strconv.Itoa(code) + `,"message":"` + err + `"}` + "\n"
}

func TestUserHandler_Create(t *testing.T) {
	type args struct {
		user *model.User
	}
	type mockBehavior func(r *mock_service.MockUser, args args)

	tests := []struct {
		name                 string
		inputBody            string
		input                args
		mock                 mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "user 1"}`,
			input: args{
				&model.User{
					Username: "user 1",
				},
			},
			mock: func(r *mock_service.MockUser, args args) {
				r.EXPECT().Create(args.user).Return("00000000-0000-0000-0000-000000000001", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":"00000000-0000-0000-0000-000000000001"}` + "\n",
		},
		{
			name:      "Bad request",
			inputBody: `{"username": ""}`,
			input: args{
				&model.User{
					Username: "",
				},
			},
			mock: func(r *mock_service.MockUser, args args) {
				r.EXPECT().Create(args.user).Return("", errMes.ErrWrongUsername)
			},
			expectedStatusCode:   400,
			expectedResponseBody: getBodyForError(400, errMes.ErrWrongUsername.Error()),
		},
		{
			name:      "Service error",
			inputBody: `{"username": "user 1"}`,
			input: args{
				&model.User{
					Username: "user 1",
				},
			},
			mock: func(r *mock_service.MockUser, args args) {
				r.EXPECT().Create(args.user).Return("", errors.New("Some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: getBodyForError(500, "Some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mock(repo, test.input)

			service := &service.Service{User: repo}
			handler := Handler{service}

			app := echo.New()
			handler.Init(app)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/add",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			app.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
