package main

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/izzanzahrial/blog-api-echo/pkg/elastic"
	"github.com/izzanzahrial/blog-api-echo/pkg/handler"
	"github.com/izzanzahrial/blog-api-echo/pkg/postgre"
	"github.com/izzanzahrial/blog-api-echo/pkg/posting"
	redisDB "github.com/izzanzahrial/blog-api-echo/pkg/redis"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"github.com/izzanzahrial/blog-api-echo/pkg/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

var (
	redisHost     = os.Getenv("redisHost")
	redisPass     = os.Getenv("redisPass")
	esAddresses   = os.Getenv("esAddresses")
	esUsername    = os.Getenv("esUsername")
	esPassword    = os.Getenv("esPassword")
	jwtSignMethod = os.Getenv("jwtSignMethod")
	jwtSignKey    = os.Getenv("jwtSignKey")
	echoAddress   = os.Getenv("echoAddress")
)

func main() {
	validator := validator.New()
	postgreDB, _ := postgre.NewPostgreDatabase()
	redis := redisDB.NewRedis(redisHost, redisPass)
	es := elastic.NewElastic(esUsername, esPassword, esAddresses)

	postRepository := repository.NewPostgre()
	postService := posting.NewService(postRepository, postgreDB, validator, redis, es)
	postHandler := handler.NewPostHandler(postService)

	userRepository := repository.NewUserPostgreRepository()
	userService := user.NewUserService(userRepository, postgreDB, validator)
	userHandler := handler.NewUserHandler(userService)

	jwtConfig := middleware.JWTConfig{
		Claims:        &user.JWTClaims{},
		SigningMethod: jwtSignMethod,
		SigningKey:    []byte(jwtSignKey),
	}

	e := echo.New()
	p := e.Group("/api/v1/posts")

	p.POST("", postHandler.Create)
	p.GET("", postHandler.FindRecent)
	p.GET("/:postid", postHandler.FindByID)
	p.PUT("/:postid", postHandler.Update, middleware.JWTWithConfig(jwtConfig))
	p.DELETE("/:postid", postHandler.Delete, middleware.JWTWithConfig(jwtConfig))
	p.GET("/:result", postHandler.FindByTitleContent)

	u := e.Group("/api/v1/user")

	u.POST("", userHandler.Create)
	u.PUT("", userHandler.UpdateUser, middleware.JWTWithConfig(jwtConfig))
	u.DELETE("", userHandler.Delete, middleware.JWTWithConfig(jwtConfig))
	u.POST("/login", userHandler.Login)
	u.PUT("/password", userHandler.UpdatePassword, middleware.JWTWithConfig(jwtConfig))

	e.Logger.Fatal(e.Start(echoAddress))
}
