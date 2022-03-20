package post

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

type PostHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	FindByID(c echo.Context) error
	FindAll(c echo.Context) error
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

func (ph *postHandler) Create(c echo.Context) error {
	// post := &entity.Post{}
	// post.Title = c.FormValue("title")
	// post.Content = c.FormValue("content")

	// postResponse, _ := ph.PostService.Create(c.Request().Context(), *post)
	// webResponse := webResponse{
	// 	code:   http.StatusAccepted,
	// 	status: "",
	// 	data:   postResponse,
	// }

	// return webResponse

	post := entity.Post{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	postResponse, err := ph.PostService.Create(c.Request().Context(), post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		code:   http.StatusCreated,
		status: "",
		data:   postResponse,
	}

	return c.JSON(http.StatusCreated, webResponse)
}

func (ph *postHandler) Update(c echo.Context) error {
	// post := &entity.Post{}
	// id := c.FormValue("id")
	// newId, err := strconv.ParseUint(id, 10, 64)
	// if err != nil {

	// }
	// post.ID = newId
	// post.Title = c.FormValue("title")
	// post.Content = c.FormValue("content")

	// postResponse, _ := ph.PostService.Update(c.Request().Context(), *post)
	// webResponse := webResponse{
	// 	code:   http.StatusAccepted,
	// 	status: "",
	// 	data:   postResponse,
	// }

	// return webResponse

	post := entity.Post{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	postResponse, err := ph.PostService.Update(c.Request().Context(), post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		code:   http.StatusAccepted,
		status: "",
		data:   postResponse,
	}

	return c.JSON(http.StatusAccepted, webResponse)
}

func (ph *postHandler) Delete(c echo.Context) error {
	// id := c.FormValue("id")
	// newId, err := strconv.ParseUint(id, 10, 64)
	// if err != nil {

	// }

	// ph.PostService.Delete(c.Request().Context(), newId)
	// webResponse := webResponse{
	// 	code:   http.StatusOK,
	// 	status: "",
	// }

	// return webResponse

	post := entity.Post{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ph.PostService.Delete(c.Request().Context(), post.ID)
	webResponse := webResponse{
		code:   http.StatusOK,
		status: "",
	}

	return c.JSON(http.StatusOK, webResponse)
}

func (ph *postHandler) FindByID(c echo.Context) error {
	// id := c.QueryParam("id")
	// newId, err := strconv.ParseUint(id, 10, 64)
	// if err != nil {

	// }

	// postResponse, _ := ph.PostService.FindByID(c.Request().Context(), newId)
	// webResponse := webResponse{
	// 	code:   http.StatusFound,
	// 	status: "",
	// 	data:   postResponse,
	// }

	// return webResponse

	post := entity.Post{}

	defer c.Request().Body.Close()

	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	postResponse, err := ph.PostService.FindByID(c.Request().Context(), post.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		code:   http.StatusFound,
		status: "",
		data:   postResponse,
	}

	return c.JSON(http.StatusFound, webResponse)
}

func (ph *postHandler) FindAll(c echo.Context) error {
	// postResponse, _ := ph.PostService.FindAll(c.Request().Context())
	// webResponse := webResponse{
	// 	code:   http.StatusOK,
	// 	status: "",
	// 	data:   postResponse,
	// }

	// return webResponse

	var posts []entity.Post

	posts, err := ph.PostService.FindAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	webResponse := webResponse{
		code:   http.StatusFound,
		status: "",
		data:   posts,
	}

	return c.JSON(http.StatusFound, webResponse)
}
