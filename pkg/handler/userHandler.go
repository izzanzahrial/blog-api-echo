package handler

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"github.com/izzanzahrial/blog-api-echo/pkg/user"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	UserService user.UserService
}

func NewUserHandler(us user.UserService) UserHandler {
	return &userHandler{
		UserService: us,
	}
}

func (us *userHandler) Create(c echo.Context) error {
	// Should i use body request or form value
	// defer c.Request().Body.Close()

	// err := json.NewDecoder(c.Request().Body).Decode(&user)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	var user user.User
	user.Name = c.FormValue("name")
	user.Password = c.FormValue("password")
	if password2 := c.FormValue("password2"); password2 != user.Password {
		return echo.ErrBadRequest
	}

	userResponse, err := us.UserService.Create(c.Request().Context(), user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	webResponse := webResponse{
		Code:    http.StatusCreated,
		Message: http.StatusText(http.StatusCreated),
		Data:    userResponse,
	}

	return c.JSON(http.StatusCreated, webResponse)
}

func (us *userHandler) UpdateUser(c echo.Context) error {
	// defer c.Request().Body.Close()

	// err := json.NewDecoder(c.Request().Body).Decode(&user)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	userClaims := c.Get("user").(*jwt.Token)
	claims := userClaims.Claims.(*user.JWTClaims)

	user := repository.User{
		ID:       claims.ID,
		Email:    claims.Email,
		Username: claims.Username,
		Name:     claims.Name,
	}

	user.Email = c.FormValue("email")
	user.Username = c.FormValue("username")
	user.Name = c.FormValue("name")

	userResponse, err := us.UserService.UpdateUser(context.Background(), user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	webResponse := webResponse{
		Code:    http.StatusAccepted,
		Message: http.StatusText(http.StatusAccepted),
		Data:    userResponse,
	}

	return c.JSON(http.StatusAccepted, webResponse)
}

func (us *userHandler) UpdatePassword(c echo.Context) error {
	// defer c.Request().Body.Close()

	// err := json.NewDecoder(c.Request().Body).Decode(&user)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	userClaims := c.Get("user").(*jwt.Token)
	claims := userClaims.Claims.(*user.JWTClaims)

	user := repository.User{
		ID:       claims.ID,
		Email:    claims.Email,
		Username: claims.Username,
		Name:     claims.Name,
	}
	user.Password = c.FormValue("password")
	if password2 := c.FormValue("password2"); password2 != user.Password {
		return echo.ErrBadRequest
	}

	newPass := c.FormValue("new_pass")

	userResponse, err := us.UserService.UpdatePassword(context.Background(), user, newPass)
	if err != nil {
		return echo.ErrInternalServerError
	}

	webResponse := webResponse{
		Code:    http.StatusAccepted,
		Message: http.StatusText(http.StatusAccepted),
		Data:    userResponse,
	}

	return c.JSON(http.StatusAccepted, webResponse)
}

func (us *userHandler) Delete(c echo.Context) error {
	// defer c.Request().Body.Close()

	// err := json.NewDecoder(c.Request().Body).Decode(&user)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	userClaims := c.Get("user").(*jwt.Token)
	claims := userClaims.Claims.(*user.JWTClaims)
	user := repository.User{
		ID:       claims.ID,
		Email:    claims.Email,
		Username: claims.Username,
		Name:     claims.Name,
	}

	user.Password = c.FormValue("password")
	if password2 := c.FormValue("password2"); password2 != user.Password {
		return echo.ErrBadRequest
	}

	if err := us.UserService.Delete(c.Request().Context(), user); err != nil {
		return echo.ErrInternalServerError
	}

	webResponse := webResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
	}

	return c.JSON(http.StatusOK, webResponse)
}

func (us *userHandler) Login(c echo.Context) error {
	// defer c.Request().Body.Close()

	// decoder := json.NewDecoder(c.Request().Body)
	// err := decoder.Decode(&user)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	emailOrUname := c.FormValue("email_or_username")
	password := c.FormValue("password")

	userResponse, token, err := us.UserService.Login(context.Background(), emailOrUname, password)
	if err != nil {
		return echo.ErrUnauthorized
	}

	// check this again
	webResponse := webResponse{
		Code:    http.StatusFound,
		Message: http.StatusText(http.StatusFound),
		Data:    []interface{}{userResponse, token},
	}

	return c.JSON(http.StatusOK, webResponse)
}
