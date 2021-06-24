package handler

import (
	"net/http"

	errMes "github.com/VolkovEgor/sellerx-task/internal/error_message"
	"github.com/VolkovEgor/sellerx-task/internal/model"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initChatRoutes(api *echo.Group) {
	chats := api.Group("/chats")
	{
		chats.POST("/add", h.createChat)
	}
}

type chatInput struct {
	Name  string `json:"name" valid:"length(1|50)"`
	Users []int  `json:"users"`
}

// @Summary Create Chat
// @Tags chats
// @Description Create chat
// @ModuleID createChat
// @Accept json
// @Produce json
// @Param input body chatInput true "chat input"
// @Success 200 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /chats/add [post]
func (h *Handler) createChat(ctx echo.Context) error {
	var input chatInput

	if err := ctx.Bind(&input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err)
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err)
	}

	chat := &model.Chat{
		Name:  input.Name,
		Users: input.Users,
	}

	chatId, err := h.services.Chat.Create(chat)
	if err != nil {
		if err == errMes.ErrWrongChatname || err == errMes.ErrNoChatUsers || err == errMes.ErrChatUserNotExists {
			return SendError(ctx, http.StatusBadRequest, err)
		}
		return SendError(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, idResponse{chatId})
}
