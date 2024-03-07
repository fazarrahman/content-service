package service

import (
	"net/http"

	"github.com/fazarrahman/content-service/domain/image/entity"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Service) UploadImage(ctx echo.Context, image *entity.Image, userID uuid.UUID) *echo.HTTPError {
	if image == nil || image.ImageBase64Str == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Image is required")
	}

	if image.Group == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Group is required")
	}

	path, err := s.storageRepository.UploadImage(ctx, image.ImageBase64Str, image.Group, userID)
	if err != nil {
		return err
	}

	image.Path = path
	return s.imageRepository.Save(ctx, image)
}
