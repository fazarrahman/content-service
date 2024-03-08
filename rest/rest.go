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
	content := e.Group("/api/v1/content", jwtlib.Required())
	content.POST("/image", r.UploadImage)
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
