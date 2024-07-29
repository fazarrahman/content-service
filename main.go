package main

import (
	"fmt"
	"net/http"

	"github.com/fazarrahman/content-service/config/postgres"
	imagePostgreRepo "github.com/fazarrahman/content-service/domain/image/repository/postgres"
	"github.com/fazarrahman/content-service/domain/storage/repository/spacebucket"
	"github.com/fazarrahman/content-service/lib/envLib"
	"github.com/fazarrahman/content-service/rest"
	"github.com/fazarrahman/content-service/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowHeaders, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodOptions, http.MethodPost, http.MethodDelete, http.MethodPut},
	}))
	db := postgres.Connection()
	imagePostgre := imagePostgreRepo.NewPostgres(db)
	spaceBucket := spacebucket.New()
	svc := service.New(imagePostgre, spaceBucket)
	rest.New(svc).Register(e)
	fmt.Println("App run at port " + envLib.GetEnv("CONTENT_APP_PORT"))
	e.Logger.Fatal(e.Start(":" + envLib.GetEnv("CONTENT_APP_PORT")))
}
