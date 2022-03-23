package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/izzanzahrial/blog-api-echo/entity"
	"github.com/labstack/echo/v4"
)

var (
	ErrFailedToDecodeBody = errors.New("decoder failed to decode the body request")
)

type UserHandler interface {
	Create(c echo.Context) error
	UpdateUser(c echo.Context) error
	UpdatePassword(c echo.Context) error
	Delete(c echo.Context) error
	Login(c echo.Context) error
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

	// Should i use body request or form value
	// defer c.Request().Body.Close()

	// err := json.NewDecoder(c.Request().Body).Decode(&user)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	user.Name = c.FormValue("name")
	user.Password = c.FormValue("password")
	if password2 := c.FormValue("password2"); password2 != user.Password {
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

func (us *userHandler) UpdateUser(c echo.Context) error {
	user := entity.User{}

	// defer c.Request().Body.Close()

	// err := json.NewDecoder(c.Request().Body).Decode(&user)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	user.Name = c.FormValue("name")

	userResponse, err := us.UserService.UpdateUser(c.Request().Context(), user)
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

func (us *userHandler) UpdatePassword(c echo.Context) error {
	user := entity.User{}

	// defer c.Request().Body.Close()

	// err := json.NewDecoder(c.Request().Body).Decode(&user)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	user.Password = c.FormValue("password")
	if password2 := c.FormValue("password2"); password2 != user.Password {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	userResponse, err := us.UserService.UpdatePassword(c.Request().Context(), user)
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

	// defer c.Request().Body.Close()

	// err := json.NewDecoder(c.Request().Body).Decode(&user)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	id := c.FormValue("id")
	id2, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	user.ID = id2
	user.Password = c.FormValue("password")

	us.UserService.Delete(c.Request().Context(), user.ID, user.Password)
	webResponse := webResponse{
		code:   http.StatusOK,
		status: "",
	}

	return c.JSON(http.StatusOK, webResponse)
}

func (us *userHandler) Login(c echo.Context) error {
	user := entity.User{}

	// defer c.Request().Body.Close()

	// decoder := json.NewDecoder(c.Request().Body)
	// err := decoder.Decode(&user)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	id := c.FormValue("id")
	id2, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}
	user.ID = id2
	user.Password = c.FormValue("password")

	userResponse, err := us.UserService.Login(c.Request().Context(), user.ID, user.Password)
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
