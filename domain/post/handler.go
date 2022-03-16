package post

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/v4"
)

type PostHandler interface {
	Create(c echo.Context)
	Update(c echo.Context)
	Delete(c echo.Context)
	FindByID(c echo.Context)
	FindAll(c echo.Context) []
}