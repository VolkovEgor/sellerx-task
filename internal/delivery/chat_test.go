package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	errMes "github.com/VolkovEgor/sellerx-task/internal/error_message"
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/service"
	mock_service "github.com/VolkovEgor/sellerx-task/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestChatHandler_Create(t *testing.T) {
	type args struct {
		chat *model.Chat
	}
	type mockBehavior func(r *mock_service.MockChat, args args)

	tests := []struct {
		name                 string
		inputBody            string
		input                args
		mock                 mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			inputBody: `{"name": "chat 1", "users": ["00000000-0000-0000-0000-000000000001", 
						"00000000-0000-0000-0000-000000000002"]}`,
			input: args{
				&model.Chat{
					Name: "chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
				},
			},
			mock: func(r *mock_service.MockChat, args args) {
				r.EXPECT().Create(args.chat).Return("00000000-0000-0000-0000-000000000001", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":"00000000-0000-0000-0000-000000000001"}` + "\n",
		},
		{
			name:      "Bad request",
			inputBody: `{"name": "chat 1", "users": []}`,
			input: args{
				&model.Chat{
					Name:  "chat 1",
					Users: []string{},
				},
			},
			mock: func(r *mock_service.MockChat, args args) {
				r.EXPECT().Create(args.chat).Return("", errMes.ErrNoChatUsers)
			},
			expectedStatusCode:   400,
			expectedResponseBody: getBodyForError(400, errMes.ErrNoChatUsers.Error()),
		},
		{
			name: "Service error",
			inputBody: `{"name": "chat 1", "users": ["00000000-0000-0000-0000-000000000001", 
						"00000000-0000-0000-0000-000000000002"]}`,
			input: args{
				&model.Chat{
					Name: "chat 1",
					Users: []string{
						"00000000-0000-0000-0000-000000000001",
						"00000000-0000-0000-0000-000000000002",
					},
				},
			},
			mock: func(r *mock_service.MockChat, args args) {
				r.EXPECT().Create(args.chat).Return("", errors.New("Some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: getBodyForError(500, "Some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockChat(c)
			test.mock(repo, test.input)

			service := &service.Service{Chat: repo}
			handler := Handler{service}

			app := echo.New()
			handler.Init(app)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/chats/add",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			app.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestChatHandler_GetAllForUser(t *testing.T) {
	type args struct {
		user string
	}
	type mockBehavior func(r *mock_service.MockChat, args args)

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
			inputBody: `{"user": "00000000-0000-0000-0000-000000000001"}`,
			input: args{
				user: "00000000-0000-0000-0000-000000000001",
			},
			mock: func(r *mock_service.MockChat, args args) {
				data := []*model.Chat{
					{
						Id:   "00000000-0000-0000-0000-000000000001",
						Name: "Chat 1",
						Users: []string{
							"00000000-0000-0000-0000-000000000001",
							"00000000-0000-0000-0000-000000000002",
						},
						CreatedAt:       100,
						LastMessageTime: 300,
					},
					{
						Id:   "00000000-0000-0000-0000-000000000002",
						Name: "Chat 2",
						Users: []string{
							"00000000-0000-0000-0000-000000000003",
							"00000000-0000-0000-0000-000000000001",
						},
						CreatedAt:       150,
						LastMessageTime: 200,
					},
				}
				r.EXPECT().GetAllForUser(args.user).Return(data, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[{"id":"00000000-0000-0000-0000-000000000001","name":"Chat 1",` +
				`"users":["00000000-0000-0000-0000-000000000001","00000000-0000-0000-0000-000000000002"],` +
				`"created_at":100,"last_message_time":300},` +
				`{"id":"00000000-0000-0000-0000-000000000002","name":"Chat 2",` +
				`"users":["00000000-0000-0000-0000-000000000003","00000000-0000-0000-0000-000000000001"],` +
				`"created_at":150,"last_message_time":200}]` + "\n",
		},
		{
			name:      "Bad request",
			inputBody: `{"user": ""}`,
			input: args{
				user: "",
			},
			mock: func(r *mock_service.MockChat, args args) {
				r.EXPECT().GetAllForUser(args.user).Return(nil, errMes.ErrEmptyUserId)
			},
			expectedStatusCode:   400,
			expectedResponseBody: getBodyForError(400, errMes.ErrEmptyUserId.Error()),
		},
		{
			name:      "Service error",
			inputBody: `{"user": "00000000-0000-0000-0000-000000000001"}`,
			input: args{
				user: "00000000-0000-0000-0000-000000000001",
			},
			mock: func(r *mock_service.MockChat, args args) {
				r.EXPECT().GetAllForUser(args.user).Return(nil, errors.New("Some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: getBodyForError(500, "Some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockChat(c)
			test.mock(repo, test.input)

			service := &service.Service{Chat: repo}
			handler := Handler{service}

			app := echo.New()
			handler.Init(app)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/chats/get",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			app.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
