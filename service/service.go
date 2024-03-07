package service

import (
	"github.com/fazarrahman/content-service/domain/image/entity"
	imageRepository "github.com/fazarrahman/content-service/domain/image/repository"
	storageRepository "github.com/fazarrahman/content-service/domain/storage/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Service struct {
	imageRepository   imageRepository.Repository
	storageRepository storageRepository.Repository
}

func New(imageRepository imageRepository.Repository, storageRepository storageRepository.Repository) *Service {
	return &Service{imageRepository, storageRepository}
}

type Repository interface {
	UploadImage(ctx echo.Context, image *entity.Image, userID uuid.UUID) error
}
