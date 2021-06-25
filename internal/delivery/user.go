package handler

import (
	"net/http"

	errMes "github.com/VolkovEgor/sellerx-task/internal/error_message"
	"github.com/VolkovEgor/sellerx-task/internal/model"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initUserRoutes(api *echo.Group) {
	users := api.Group("/users")
	{
		users.POST("/add", h.createUser)
	}
}

type userInput struct {
	Username string `json:"username" valid:"length(1|50)"`
}

// @Summary Create User
// @Tags users
// @Description Create user
// @ModuleID createUser
// @Accept json
// @Produce json
// @Param input body userInput true "user input"
// @Success 200 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/add [post]
func (h *Handler) createUser(ctx echo.Context) error {
	var input userInput

	if err := ctx.Bind(&input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err)
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return SendError(ctx, http.StatusBadRequest, err)
	}

	user := &model.User{
		Username: input.Username,
	}

	userId, err := h.services.User.Create(user)
	if err != nil {
		if err == errMes.ErrWrongUsername {
			return SendError(ctx, http.StatusBadRequest, err)
		}
		return SendError(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, idResponse{userId})
}
