package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/izzanzahrial/blog-api-echo/pkg/posting"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
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
	FindAll(c echo.Context) error
}

type postHandler struct {
	Service posting.Service
}

func NewPostHandler(ps posting.Service) PostHandler {
	return &postHandler{
		Service: ps,
	}
}

type webResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func (ph *postHandler) Create(c echo.Context) error {
	post := repository.Post{}
	post.Title = c.FormValue("title")
	post.Content = c.FormValue("content")

	// defer c.Request().Body.Close()

	// err := json.NewDecoder(c.Request().Body).Decode(&post)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	ctx := context.Background()

	postResponse, err := ph.Service.Create(ctx, post)
	if err != nil {
		return echo.ErrInternalServerError
	}

	webResponse := webResponse{
		Code:   http.StatusCreated,
		Status: "",
		Data:   postResponse,
	}

	return c.JSON(http.StatusCreated, webResponse)
}

func (ph *postHandler) Update(c echo.Context) error {
	id := c.FormValue("id")
	updatedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return echo.ErrInternalServerError
	}

	post := repository.Post{}
	post.ID = updatedID
	post.Title = c.FormValue("title")
	post.Content = c.FormValue("content")

	ctx := context.Background()

	postResponse, err := ph.Service.Update(ctx, post)
	if err != nil {
		return echo.ErrInternalServerError
	}

	webResponse := webResponse{
		Code:   http.StatusAccepted,
		Status: "",
		Data:   postResponse,
	}

	return c.JSON(http.StatusAccepted, webResponse)
}

func (ph *postHandler) Delete(c echo.Context) error {
	id := c.FormValue("id")
	deletedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return echo.ErrInternalServerError
	}

	ctx := context.Background()

	if err := ph.Service.Delete(ctx, deletedID); err != nil {
		return echo.ErrInternalServerError
	}

	webResponse := webResponse{
		Code:   http.StatusOK,
		Status: "",
	}

	return c.JSON(http.StatusOK, webResponse)
}

func (ph *postHandler) FindByID(c echo.Context) error {
	id := c.Param("id")
	searchID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return echo.ErrInternalServerError
	}

	ctx := context.Background()

	postResponse, err := ph.Service.FindByID(ctx, searchID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		Code:   http.StatusFound,
		Status: "",
		Data:   postResponse,
	}

	return c.JSON(http.StatusFound, webResponse)
}

func (ph *postHandler) FindAll(c echo.Context) error {
	var posts []repository.Post

	ctx := context.Background()

	posts, err := ph.Service.FindAll(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		Code:   http.StatusFound,
		Status: "",
		Data:   posts,
	}

	return c.JSON(http.StatusFound, webResponse)
}
