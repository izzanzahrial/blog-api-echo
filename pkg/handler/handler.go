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
	FindByTitleContent(c echo.Context) error
	FindRecent(c echo.Context) error
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
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (ph *postHandler) Create(c echo.Context) error {
	var post posting.PostData
	post.Title = c.FormValue("title")
	post.ShortDesc = c.FormValue("short_desc")
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
		Code:    http.StatusCreated,
		Message: http.StatusText(http.StatusCreated),
		Data:    postResponse,
	}

	return c.JSON(http.StatusCreated, webResponse)
}

func (ph *postHandler) Update(c echo.Context) error {
	strID := c.FormValue("id")
	updatedID, err := strconv.Atoi(strID)
	if err != nil {
		return echo.ErrInternalServerError
	}

	var post repository.PostData
	post.ID = int64(updatedID)
	post.Title = c.FormValue("title")
	post.ShortDesc = c.FormValue("short_desc")
	post.Content = c.FormValue("content")

	ctx := context.Background()

	if err = ph.Service.Update(ctx, post); err != nil {
		return echo.ErrInternalServerError
	}

	webResponse := webResponse{
		Code:    http.StatusAccepted,
		Message: http.StatusText(http.StatusAccepted),
		Data:    post,
	}

	return c.JSON(http.StatusAccepted, webResponse)
}

func (ph *postHandler) Delete(c echo.Context) error {
	strID := c.FormValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return echo.ErrInternalServerError
	}

	ctx := context.Background()

	if err := ph.Service.Delete(ctx, int64(id)); err != nil {
		return echo.ErrInternalServerError
	}

	webResponse := webResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
	}

	return c.JSON(http.StatusOK, webResponse)
}

func (ph *postHandler) FindByID(c echo.Context) error {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return echo.ErrInternalServerError
	}

	ctx := context.Background()

	postResponse, err := ph.Service.FindByID(ctx, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		Code:    http.StatusFound,
		Message: http.StatusText(http.StatusFound),
		Data:    postResponse,
	}

	return c.JSON(http.StatusFound, webResponse)
}

func (ph *postHandler) FindByTitleContent(c echo.Context) error {
	var posts []repository.PostData

	ctx := context.Background()

	query := c.QueryParam("query")
	from, err := strconv.Atoi(c.QueryParam("from"))
	if err != nil {
		return echo.ErrInternalServerError
	}

	size, err := strconv.Atoi(c.QueryParam("size"))
	if err != nil {
		return echo.ErrInternalServerError
	}

	posts, err = ph.Service.FindByTitleContent(ctx, query, from, size)
	if err != nil {
		return echo.ErrInternalServerError
	}

	webResponse := webResponse{
		Code:    http.StatusFound,
		Message: http.StatusText(http.StatusFound),
		Data:    posts,
	}

	return c.JSON(http.StatusOK, webResponse)
}

func (ph *postHandler) FindRecent(c echo.Context) error {
	var posts []repository.PostData

	ctx := context.Background()

	from, err := strconv.Atoi(c.QueryParam("from"))
	if err != nil {
		return echo.ErrInternalServerError
	}

	size, err := strconv.Atoi(c.QueryParam("size"))
	if err != nil {
		return echo.ErrInternalServerError
	}

	posts, err = ph.Service.FindRecent(ctx, from, size)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		Code:    http.StatusFound,
		Message: http.StatusText(http.StatusFound),
		Data:    posts,
	}

	return c.JSON(http.StatusFound, webResponse)
}
