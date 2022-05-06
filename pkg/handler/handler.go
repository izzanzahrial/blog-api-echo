package handler

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var (
	ErrFailedToDecodeBody = errors.New("decoder failed to decode the body request")
)

type PostHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	FindByID(c echo.Context) error
	FindByTitleContent(c echo.Context) error
	FindRecent(c echo.Context) error
}

type UserHandler interface {
	Create(c echo.Context) error
	UpdateUser(c echo.Context) error
	UpdatePassword(c echo.Context) error
	Delete(c echo.Context) error
	Login(c echo.Context) error
}

type webResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
