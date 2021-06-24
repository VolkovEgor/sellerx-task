package handler

import (
	"net/http"

	_ "github.com/VolkovEgor/sellerx-task/docs/swagger"
	"github.com/VolkovEgor/sellerx-task/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type errorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

type idResponse struct {
	Id int `json:"id"`
}

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

// @title SellerX Task API
// @version 1.0
// @description API Server for SellerX Task

// @host localhost:9000
// @BasePath /

func (h *Handler) Init(router *echo.Echo) {
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	router.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api := router.Group("")
	{
		h.initUserRoutes(api)
		h.initChatRoutes(api)
	}
}

func SendError(ctx echo.Context, status int, err error) error {
	logrus.Error(err.Error())
	return ctx.JSON(status, errorResponse{status, err.Error()})
}
