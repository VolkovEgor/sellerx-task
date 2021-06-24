package handler

import (
	"net/http"

	errMes "github.com/VolkovEgor/sellerx-task/internal/error_message"
	"github.com/VolkovEgor/sellerx-task/internal/model"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initMessageRoutes(api *echo.Group) {
	messages := api.Group("/messages")
	{
		messages.POST("/add", h.createMessage)
	}
}

type messageInput struct {
	ChatId   int    `json:"chat"`
	AuthorId int    `json:"author"`
	Text     string `json:"text"`
}

// @Summary Create Message
// @Tags messages
// @Description Create message
// @ModuleID createMessage
// @Accept json
// @Produce json
// @Param input body messageInput true "message input"
// @Success 200 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /messages/add [post]
func (h *Handler) createMessage(ctx echo.Context) error {
	var input messageInput

	if err := ctx.Bind(&input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err)
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err)
	}

	message := &model.Message{
		ChatId:   input.ChatId,
		AuthorId: input.AuthorId,
		Text:     input.Text,
	}

	messageId, err := h.services.Message.Create(message)
	if err != nil {
		if err == errMes.ErrWrongTextMes || err == errMes.ErrAuthorNotExists || err == errMes.ErrChatNotExists {
			return SendError(ctx, http.StatusBadRequest, err)
		}
		return SendError(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, idResponse{messageId})
}
