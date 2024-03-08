package minio

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type MinioStorage struct {
	client   *minio.Client
	bucket   string
	endpoint string
}

func New() *MinioStorage {
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	secretKey := os.Getenv("STORAGE_SECRET_KEY")
	bucketName := os.Getenv("STORAGE_BUCKET_NAME")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		logrus.Errorln("Failed to connect to MinIO:", err)
	}

	return &MinioStorage{
		client:   minioClient,
		bucket:   bucketName,
		endpoint: endpoint,
	}
}

func (s *MinioStorage) UploadImage(ctx echo.Context, imageBase64Str, groupName string, userID uuid.UUID) (string, *echo.HTTPError) {
	b64data := imageBase64Str[strings.IndexByte(imageBase64Str, ',')+1:]
	decodedImage, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to decode image: %s", err.Error()))
	}

	fileType := http.DetectContentType(decodedImage)
	validFileExtensions := []string{".png", ".jpg", ".jpeg"}
	var fileExtension string

	for _, ext := range validFileExtensions {
		if mime.TypeByExtension(ext) == fileType {
			fileExtension = ext
			break
		}
	}

	if fileExtension == "" {
		return "", echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("file extension must be : %v", validFileExtensions))
	}

	if len(decodedImage) > 2*1024*1024 {
		return "", echo.NewHTTPError(http.StatusBadRequest, "image is too large (2MB max)")
	}

	filename := uuid.New().String() + fileExtension
	objectName := fmt.Sprintf("%s/%s/%s", userID.String(), groupName, filename)

	info, err := s.client.PutObject(context.Background(), s.bucket, objectName, bytes.NewReader(decodedImage), int64(len(decodedImage)), minio.PutObjectOptions{})
	if err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to upload image: %s", err.Error()))
	}

	return fmt.Sprintf("%s/%s/%s", os.Getenv("STORAGE_PUBLIC_URL"), info.Bucket, info.Key), nil
}
