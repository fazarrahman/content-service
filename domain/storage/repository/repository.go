package repository

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	UploadImage(ctx echo.Context, File, groupName string, userID uuid.UUID) (string, *echo.HTTPError)
}
