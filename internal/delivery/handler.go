package handler

import (
	"net/http"

	"github.com/VolkovEgor/sellerx-task/internal/service"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

// @title SellerX Task
// @version 1.0
// @description API Server for SellerX Task

// @host localhost:9000
// @BasePath /api/

func (h *Handler) Init(router *echo.Echo) {
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	router.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
}
