package service

import (
	"net/http"
	"strings"

	"github.com/fazarrahman/content-service/domain/image/entity"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Service) UploadImage(ctx echo.Context, image *entity.Image, userID uuid.UUID) *echo.HTTPError {
	if image == nil || image.File == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Image is required")
	} else if image.Group == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Group is required")
	} else if strings.TrimSpace(image.Title) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Title is required")
	} else if strings.TrimSpace(image.Description) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Description is required")
	} else if len(image.TagStr) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Image must have tag at least 1")
	}

	image.UserID = userID
	path, err := s.storageRepository.UploadImage(ctx, image, userID)
	if err != nil {
		return err
	}

	image.Path = path
	image.CreatedBy = userID
	image.UpdatedBy = userID
	return s.imageRepository.Save(ctx, image)
}

func (s *Service) GetList(ctx echo.Context) ([]*entity.Image, *echo.HTTPError) {
	return s.imageRepository.GetList(ctx)
}
