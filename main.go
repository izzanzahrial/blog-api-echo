package main

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/izzanzahrial/blog-api-echo/domain/post"
	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"
)

type JWTCustomClaims struct {
	Name  string
	Admin bool
	jwt.StandardClaims
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "izzan" || password != "blablabla" {
		return echo.ErrUnauthorized
	}

	claims := &JWTCustomClaims{
		"izzan",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func main() {
	validate := validator.New()
	postgreDB, _ := post.NewPostgreDatabase()

	postRepository := post.NewPostgreRepository()
	postService := post.NewPostService(postRepository, postgreDB, validate)
	postHandler := post.NewPostHandler(postService)

	e := echo.New()
	p := e.Group("/api/v1/posts")

	p.GET("", postHandler.FindAll)
	p.GET("/:postid", postHandler.FindByID)
	p.POST("", postHandler.Create)
	p.PUT("/:postid", postHandler.Update)
	p.DELETE("/:postid", postHandler.Delete)

	e.Logger.Fatal(e.Start(":1323"))
}
