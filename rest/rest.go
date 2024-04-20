package rest

import (
	"net/http"

	"github.com/fazarrahman/content-service/domain/image/entity"
	jwtlib "github.com/fazarrahman/content-service/lib/jwtLib"
	"github.com/fazarrahman/content-service/service"
	"github.com/labstack/echo/v4"
)

type Rest struct {
	service *service.Service
}

func New(service *service.Service) *Rest {
	return &Rest{service: service}
}

func (r *Rest) Register(e *echo.Echo) {
	content := e.Group("/api/v1/content")
	content.POST("/image", r.UploadImage, jwtlib.Required())
	content.GET("/image", r.GetList)
}

func (r *Rest) UploadImage(c echo.Context) error {
	image := entity.Image{}
	if err := c.Bind(&image); err != nil {
		return c.JSON(http.StatusBadRequest, "Error parsing request body")
	}
	claims := jwtlib.GetClaims(c)
	err := r.service.UploadImage(c, &image, claims.UserId)
	if err != nil {
		return c.JSON(err.Code, echo.Map{"message": err.Message})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Success"})
}

func (r *Rest) GetList(c echo.Context) error {
	images, err := r.service.GetList(c)
	if err != nil {
		return c.JSON(err.Code, echo.Map{"message": err.Message})
	}
	return c.JSON(http.StatusOK, echo.Map{"data": images})
}
