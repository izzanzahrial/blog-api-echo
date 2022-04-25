package main

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/izzanzahrial/blog-api-echo/domain/post"
	"github.com/izzanzahrial/blog-api-echo/domain/user"
	"github.com/izzanzahrial/blog-api-echo/pkg/elastic"
	"github.com/izzanzahrial/blog-api-echo/pkg/handler"
	"github.com/izzanzahrial/blog-api-echo/pkg/posting"
	"github.com/izzanzahrial/blog-api-echo/pkg/redis"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

var (
	redisHost   = os.Getenv("redisHost")
	redisPass   = os.Getenv("redisPass")
	esAddresses = os.Getenv("esAddresses")
	esUsername  = os.Getenv("esUsername")
	esPassword  = os.Getenv("esPassword")
)

func main() {
	validator := validator.New()
	postgreDB, _ := post.NewPostgreDatabase()
	redis := redis.NewRedis(redisHost, redisPass)
	es := elastic.NewElastic(esUsername, esPassword, esAddresses)

	postRepository := repository.NewPostgre()
	postService := posting.NewService(postRepository, postgreDB, validator, redis, es.Client)
	postHandler := handler.NewPostHandler(postService)

	userRepository := user.NewPostgreRepository()
	userService := user.NewUserService(userRepository, postgreDB, validator)
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
