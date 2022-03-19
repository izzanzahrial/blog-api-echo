package post

import (
	"net/http"
	"strconv"

	"github.com/izzanzahrial/blog-api-echo/entity"
	"github.com/labstack/echo/v4"
)

type PostHandler interface {
	Create(c echo.Context) webResponse
	Update(c echo.Context) webResponse
	Delete(c echo.Context) webResponse
	FindByID(c echo.Context) webResponse
	FindAll(c echo.Context) webResponse
}

type postHandler struct {
	PostService PostService
}

func NewPostHandler(ps PostService) PostHandler {
	return &postHandler{
		PostService: ps,
	}
}

type webResponse struct {
	code   int
	status string
	data   interface{}
}

func (ph *postHandler) Create(c echo.Context) webResponse {
	post := &entity.Post{}
	post.Title = c.FormValue("title")
	post.Content = c.FormValue("content")

	postResponse, _ := ph.PostService.Create(c.Request().Context(), *post)
	webResponse := webResponse{
		code:   http.StatusAccepted,
		status: "",
		data:   postResponse,
	}

	return webResponse
}

func (ph *postHandler) Update(c echo.Context) webResponse {
	post := &entity.Post{}
	id := c.FormValue("id")
	newId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {

	}
	post.ID = newId
	post.Title = c.FormValue("title")
	post.Content = c.FormValue("content")

	postResponse, _ := ph.PostService.Update(c.Request().Context(), *post)
	webResponse := webResponse{
		code:   http.StatusAccepted,
		status: "",
		data:   postResponse,
	}

	return webResponse
}

func (ph *postHandler) Delete(c echo.Context) webResponse {
	id := c.FormValue("id")
	newId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {

	}

	ph.PostService.Delete(c.Request().Context(), newId)
	webResponse := webResponse{
		code:   http.StatusOK,
		status: "",
	}

	return webResponse
}

func (ph *postHandler) FindByID(c echo.Context) webResponse {
	id := c.QueryParam("id")
	newId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {

	}

	postResponse, _ := ph.PostService.FindByID(c.Request().Context(), newId)
	webResponse := webResponse{
		code:   http.StatusFound,
		status: "",
		data:   postResponse,
	}

	return webResponse
}

func (ph *postHandler) FindAll(c echo.Context) webResponse {
	postResponse, _ := ph.PostService.FindAll(c.Request().Context())
	webResponse := webResponse{
		code:   http.StatusOK,
		status: "",
		data:   postResponse,
	}

	return webResponse
}
