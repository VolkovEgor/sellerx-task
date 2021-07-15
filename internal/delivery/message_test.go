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

func TestMessageHandler_Create(t *testing.T) {
	type args struct {
		message *model.Message
	}
	type mockBehavior func(r *mock_service.MockMessage, args args)

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
			inputBody: `{"chat": "00000000-0000-0000-0000-000000000001", 
						"author": "00000000-0000-0000-0000-000000000001", "text": "Message 1"}`,
			input: args{
				message: &model.Message{
					ChatId:   "00000000-0000-0000-0000-000000000001",
					AuthorId: "00000000-0000-0000-0000-000000000001",
					Text:     "Message 1",
				},
			},
			mock: func(r *mock_service.MockMessage, args args) {
				r.EXPECT().Create(args.message).Return("00000000-0000-0000-0000-000000000001", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":"00000000-0000-0000-0000-000000000001"}` + "\n",
		},
		{
			name:      "Bad request",
			inputBody: `{"chat": "00000000-0000-0000-0000-000000000001", "author": "", "text": "message 1"}`,
			input: args{
				message: &model.Message{
					ChatId:   "00000000-0000-0000-0000-000000000001",
					AuthorId: "",
					Text:     "message 1",
				},
			},
			mock: func(r *mock_service.MockMessage, args args) {
				r.EXPECT().Create(args.message).Return("", errMes.ErrEmptyUserId)
			},
			expectedStatusCode:   400,
			expectedResponseBody: getBodyForError(400, errMes.ErrEmptyUserId.Error()),
		},
		{
			name: "Service error",
			inputBody: `{"chat": "00000000-0000-0000-0000-000000000001", 
						"author": "00000000-0000-0000-0000-000000000001", "text": "Message 1"}`,
			input: args{
				message: &model.Message{
					ChatId:   "00000000-0000-0000-0000-000000000001",
					AuthorId: "00000000-0000-0000-0000-000000000001",
					Text:     "Message 1",
				},
			},
			mock: func(r *mock_service.MockMessage, args args) {
				r.EXPECT().Create(args.message).Return("", errors.New("Some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: getBodyForError(500, "Some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockMessage(c)
			test.mock(repo, test.input)

			service := &service.Service{Message: repo}
			handler := Handler{service}

			app := echo.New()
			handler.Init(app)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/messages/add",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			app.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestMessageHandler_GetAllForChat(t *testing.T) {
	type args struct {
		chat string
	}
	type mockBehavior func(r *mock_service.MockMessage, args args)

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
			inputBody: `{"chat": "00000000-0000-0000-0000-000000000001"}`,
			input: args{
				chat: "00000000-0000-0000-0000-000000000001",
			},
			mock: func(r *mock_service.MockMessage, args args) {
				data := []*model.Message{
					{
						Id:        "00000000-0000-0000-0000-000000000001",
						ChatId:    "00000000-0000-0000-0000-000000000001",
						AuthorId:  "00000000-0000-0000-0000-000000000001",
						Text:      "Message 1",
						CreatedAt: 100,
					},
					{
						Id:        "00000000-0000-0000-0000-000000000002",
						ChatId:    "00000000-0000-0000-0000-000000000001",
						AuthorId:  "00000000-0000-0000-0000-000000000002",
						Text:      "Message 2",
						CreatedAt: 200,
					},
				}
				r.EXPECT().GetAllForChat(args.chat).Return(data, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[{"id":"00000000-0000-0000-0000-000000000001",` +
				`"chat_id":"00000000-0000-0000-0000-000000000001","author_id":"00000000-0000-0000-0000-000000000001",` +
				`"text":"Message 1","created_at":100},` +
				`{"id":"00000000-0000-0000-0000-000000000002",` +
				`"chat_id":"00000000-0000-0000-0000-000000000001","author_id":"00000000-0000-0000-0000-000000000002",` +
				`"text":"Message 2","created_at":200}]` + "\n",
		},
		{
			name:      "Bad request",
			inputBody: `{"chat": ""}`,
			input: args{
				chat: "",
			},
			mock: func(r *mock_service.MockMessage, args args) {
				r.EXPECT().GetAllForChat(args.chat).Return(nil, errMes.ErrEmptyChatId)
			},
			expectedStatusCode:   400,
			expectedResponseBody: getBodyForError(400, errMes.ErrEmptyChatId.Error()),
		},
		{
			name:      "Service error",
			inputBody: `{"chat": "00000000-0000-0000-0000-000000000001"}`,
			input: args{
				chat: "00000000-0000-0000-0000-000000000001",
			},
			mock: func(r *mock_service.MockMessage, args args) {
				r.EXPECT().GetAllForChat(args.chat).Return(nil, errors.New("Some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: getBodyForError(500, "Some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockMessage(c)
			test.mock(repo, test.input)

			service := &service.Service{Message: repo}
			handler := Handler{service}

			app := echo.New()
			handler.Init(app)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/messages/get",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			app.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
