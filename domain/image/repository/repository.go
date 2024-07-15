package repository

import (
	"github.com/fazarrahman/content-service/domain/image/entity"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	Save(ctx echo.Context, image *entity.Image) *echo.HTTPError
	GetList(ctx echo.Context, page, size int) ([]*entity.Image, *echo.HTTPError)
}
