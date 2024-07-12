package repository

import (
	"github.com/fazarrahman/content-service/domain/image/entity"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	UploadImage(ctx echo.Context, image *entity.Image, userID uuid.UUID) (string, *echo.HTTPError)
}
