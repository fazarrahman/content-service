package repository

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	UploadImage(ctx echo.Context, imageBase64Str, groupName string, userID uuid.UUID) (string, *echo.HTTPError)
}
