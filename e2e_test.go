package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	handler "github.com/VolkovEgor/sellerx-task/internal/delivery"
	"github.com/VolkovEgor/sellerx-task/internal/model"
	"github.com/VolkovEgor/sellerx-task/internal/repository"
	"github.com/VolkovEgor/sellerx-task/internal/service"
	"github.com/VolkovEgor/sellerx-task/test"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
)

type idResponse struct {
	Id string `json:"id"`
}

func Test_E2E_App(t *testing.T) {
	const prefix = "./"
	db, err := test.PrepareTestDatabase(prefix, false)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handler := handler.NewHandler(service)

	app := echo.New()
	handler.Init(app)

	var userId1, userId2, userId3 idResponse
	var chatId1, chatId2 idResponse
	var messageId11, messageId12, messageId21 idResponse

	// Create first user
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
			inputUsername  = `"First Uer"`
		)
		inputBody := fmt.Sprintf(`{"username": %s}`, inputUsername)

		Convey("When create user method", func() {
			req := httptest.NewRequest(
				"POST",
				"/users/add",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if err = json.Unmarshal(w.Body.Bytes(), &userId1); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)
				So(userId1.Id, ShouldNotEqual, "")
			})
		})
	})

	// Create second user
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
			inputUsername  = `"First Advert"`
		)
		inputBody := fmt.Sprintf(`{"username": %s}`, inputUsername)

		Convey("When create user method", func() {
			req := httptest.NewRequest(
				"POST",
				"/users/add",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if err = json.Unmarshal(w.Body.Bytes(), &userId2); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)
				So(userId2.Id, ShouldNotEqual, "")
			})
		})
	})

	// Create third user
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
			inputUsername  = `"Third Advert"`
		)
		inputBody := fmt.Sprintf(`{"username": %s}`, inputUsername)

		Convey("When create user method", func() {
			req := httptest.NewRequest(
				"POST",
				"/users/add",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if err = json.Unmarshal(w.Body.Bytes(), &userId3); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)
				So(userId3.Id, ShouldNotEqual, "")
			})
		})
	})

	// Create first chat
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
			inputName      = `"First chat"`
		)
		inputUserId1 := fmt.Sprintf(`"%s"`, userId1.Id)
		inputUserId2 := fmt.Sprintf(`"%s"`, userId2.Id)
		inputUsers := []string{inputUserId1, inputUserId2}

		users := `[` + strings.Join(inputUsers, `, `) + `]`
		inputBody := fmt.Sprintf(`{"name": %s, "users": %s}`, inputName, users)

		Convey("When create chat method", func() {
			req := httptest.NewRequest(
				"POST",
				"/chats/add",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if err = json.Unmarshal(w.Body.Bytes(), &chatId1); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)
				So(chatId1.Id, ShouldNotEqual, "")
			})
		})
	})

	time.Sleep(1000 * time.Millisecond)

	// Create second chat
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
			inputName      = `"Second chat"`
		)
		inputUserId1 := fmt.Sprintf(`"%s"`, userId1.Id)
		inputUserId2 := fmt.Sprintf(`"%s"`, userId2.Id)
		inputUserId3 := fmt.Sprintf(`"%s"`, userId3.Id)
		inputUsers := []string{inputUserId1, inputUserId2, inputUserId3}

		users := `[` + strings.Join(inputUsers, `, `) + `]`
		inputBody := fmt.Sprintf(`{"name": %s, "users": %s}`, inputName, users)

		Convey("When create chat method", func() {
			req := httptest.NewRequest(
				"POST",
				"/chats/add",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if err = json.Unmarshal(w.Body.Bytes(), &chatId2); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)
				So(chatId2.Id, ShouldNotEqual, "")
			})
		})
	})
	time.Sleep(1000 * time.Millisecond)

	// Create first message in first chat
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
			inputText      = `"First message in first chat"`
		)
		inputChatId := fmt.Sprintf(`"%s"`, chatId1.Id)
		inputAuthorId := fmt.Sprintf(`"%s"`, userId2.Id)

		inputBody := fmt.Sprintf(`{"chat": %s, "author": %s, "text": %s}`,
			inputChatId, inputAuthorId, inputText)

		Convey("When create message method", func() {
			req := httptest.NewRequest(
				"POST",
				"/messages/add",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if err = json.Unmarshal(w.Body.Bytes(), &messageId11); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)
				So(messageId11.Id, ShouldNotEqual, "")
			})
		})
	})
	time.Sleep(1000 * time.Millisecond)

	// Create second message in first chat
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
			inputText      = `"Second message in first chat"`
		)
		inputChatId := fmt.Sprintf(`"%s"`, chatId1.Id)
		inputAuthorId := fmt.Sprintf(`"%s"`, userId1.Id)

		inputBody := fmt.Sprintf(`{"chat": %s, "author": %s, "text": %s}`,
			inputChatId, inputAuthorId, inputText)

		Convey("When create message method", func() {
			req := httptest.NewRequest(
				"POST",
				"/messages/add",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if err = json.Unmarshal(w.Body.Bytes(), &messageId12); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)
				So(messageId12.Id, ShouldNotEqual, "")
			})
		})
	})
	time.Sleep(1000 * time.Millisecond)

	// Create first message in second chat
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
			inputText      = `"First message in second chat"`
		)
		inputChatId := fmt.Sprintf(`"%s"`, chatId2.Id)
		inputAuthorId := fmt.Sprintf(`"%s"`, userId3.Id)

		inputBody := fmt.Sprintf(`{"chat": %s, "author": %s, "text": %s}`,
			inputChatId, inputAuthorId, inputText)

		Convey("When create message method", func() {
			req := httptest.NewRequest(
				"POST",
				"/messages/add",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if err = json.Unmarshal(w.Body.Bytes(), &messageId21); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)
				So(messageId21.Id, ShouldNotEqual, "")
			})
		})
	})

	// Get all chats for first user
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
		)
		expectedChats := []model.Chat{
			{
				Id:    chatId2.Id,
				Name:  "Second chat",
				Users: []string{userId1.Id, userId2.Id, userId3.Id},
			},
			{
				Id:    chatId1.Id,
				Name:  "First chat",
				Users: []string{userId1.Id, userId2.Id},
			},
		}

		inputUserId := fmt.Sprintf(`"%s"`, userId1.Id)
		inputBody := fmt.Sprintf(`{"user": %s}`, inputUserId)

		Convey("When get all chats for user method", func() {
			req := httptest.NewRequest(
				"POST",
				"/chats/get",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			var chats []*model.Chat
			if err = json.Unmarshal(w.Body.Bytes(), &chats); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)

				for i := 0; i < len(chats); i++ {
					So(chats[i].Id, ShouldEqual, expectedChats[i].Id)
					So(chats[i].Name, ShouldEqual, expectedChats[i].Name)
					So(chats[i].Users, ShouldResemble, expectedChats[i].Users)
				}
			})
		})
	})

	// Get all mesages for first chat
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
		)
		expectedChats := []model.Message{
			{
				Id:       messageId11.Id,
				ChatId:   chatId1.Id,
				AuthorId: userId2.Id,
				Text:     "First message in first chat",
			},
			{
				Id:       messageId12.Id,
				ChatId:   chatId1.Id,
				AuthorId: userId1.Id,
				Text:     "Second message in first chat",
			},
		}

		inputChatId := fmt.Sprintf(`"%s"`, chatId1.Id)
		inputBody := fmt.Sprintf(`{"chat": %s}`, inputChatId)

		Convey("When get all messages for chat method", func() {
			req := httptest.NewRequest(
				"POST",
				"/messages/get",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			var chats []*model.Message
			if err = json.Unmarshal(w.Body.Bytes(), &chats); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)

				for i := 0; i < len(chats); i++ {
					So(chats[i].Id, ShouldEqual, expectedChats[i].Id)
					So(chats[i].ChatId, ShouldEqual, expectedChats[i].ChatId)
					So(chats[i].AuthorId, ShouldResemble, expectedChats[i].AuthorId)
					So(chats[i].Text, ShouldResemble, expectedChats[i].Text)
				}
			})
		})
	})

	// Get all mesages for second chat
	Convey("Given params", t, func() {
		const (
			expectedStatus = http.StatusOK
		)
		expectedChats := []model.Message{
			{
				Id:       messageId21.Id,
				ChatId:   chatId2.Id,
				AuthorId: userId3.Id,
				Text:     "First message in second chat",
			},
		}

		inputChatId := fmt.Sprintf(`"%s"`, chatId2.Id)
		inputBody := fmt.Sprintf(`{"chat": %s}`, inputChatId)

		Convey("When get all messages for chat method", func() {
			req := httptest.NewRequest(
				"POST",
				"/messages/get",
				bytes.NewBufferString(inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			var chats []*model.Message
			if err = json.Unmarshal(w.Body.Bytes(), &chats); err != nil {
				log.Fatalf(err.Error())
			}

			Convey("Then should be Ok", func() {
				So(w.Code, ShouldEqual, expectedStatus)

				for i := 0; i < len(chats); i++ {
					So(chats[i].Id, ShouldEqual, expectedChats[i].Id)
					So(chats[i].ChatId, ShouldEqual, expectedChats[i].ChatId)
					So(chats[i].AuthorId, ShouldResemble, expectedChats[i].AuthorId)
					So(chats[i].Text, ShouldResemble, expectedChats[i].Text)
				}
			})
		})
	})
}
