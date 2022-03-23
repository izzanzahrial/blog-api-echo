package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/izzanzahrial/blog-api-echo/entity"
	"github.com/labstack/echo/v4"
)

var (
	ErrFailedToDecodeBody = errors.New("decoder failed to decode the body request")
)

type UserHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Find(c echo.Context) error
}

type userHandler struct {
	UserService UserService
}

func NewUserHandler(us UserService) UserHandler {
	return &userHandler{
		UserService: us,
	}
}

type webResponse struct {
	code   int
	status string
	data   interface{}
}

func (us *userHandler) Create(c echo.Context) error {
	user := entity.User{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	userResponse, err := us.UserService.Create(c.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		code:   http.StatusCreated,
		status: "",
		data:   userResponse,
	}

	return c.JSON(http.StatusCreated, webResponse)
}

func (us *userHandler) Update(c echo.Context) error {
	user := entity.User{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	userResponse, err := us.UserService.Update(c.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		code:   http.StatusAccepted,
		status: "",
		data:   userResponse,
	}

	return c.JSON(http.StatusAccepted, webResponse)
}

func (us *userHandler) Delete(c echo.Context) error {
	user := entity.User{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	us.UserService.Delete(c.Request().Context(), user.ID, user.Password)
	webResponse := webResponse{
		code:   http.StatusOK,
		status: "",
	}

	return c.JSON(http.StatusOK, webResponse)
}

func (us *userHandler) Find(c echo.Context) error {
	user := entity.User{}

	defer c.Request().Body.Close()

	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	userResponse, err := us.UserService.Find(c.Request().Context(), user.ID, user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		code:   http.StatusFound,
		status: "",
		data:   userResponse,
	}

	return c.JSON(http.StatusFound, webResponse)
}
