package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/izzanzahrial/blog-api-echo/domain/post"
	"github.com/izzanzahrial/blog-api-echo/domain/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

func main() {
	validate := validator.New()
	postgreDB, _ := post.NewPostgreDatabase()

	postRepository := post.NewPostgreRepository()
	postService := post.NewPostService(postRepository, postgreDB, validate)
	postHandler := post.NewPostHandler(postService)

	userRepository := user.NewPostgreRepository()
	userService := user.NewUserService(userRepository, postgreDB, validate)
	userHandler := user.NewUserHandler(userService)

	jwtConfig := middleware.JWTConfig{
		Claims:        &user.JWTClaims{},
		SigningMethod: "HS512",
		SigningKey:    []byte("izzan"),
	}

	e := echo.New()
	p := e.Group("/api/v1/posts")

	p.POST("", postHandler.Create)
	p.GET("", postHandler.FindAll)
	p.GET("/:postid", postHandler.FindByID)
	p.PUT("/:postid", postHandler.Update, middleware.JWTWithConfig(jwtConfig))
	p.DELETE("/:postid", postHandler.Delete, middleware.JWTWithConfig(jwtConfig))

	u := e.Group("/api/v1/user")

	u.POST("", userHandler.Create)
	u.PUT("", userHandler.UpdateUser, middleware.JWTWithConfig(jwtConfig))
	u.DELETE("", userHandler.Delete, middleware.JWTWithConfig(jwtConfig))
	u.POST("/login", userHandler.Login)
	u.PUT("/password", userHandler.UpdatePassword, middleware.JWTWithConfig(jwtConfig))

	e.Logger.Fatal(e.Start(":1323"))
}
