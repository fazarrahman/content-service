package spacebucket

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/fazarrahman/content-service/domain/image/entity"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SpaceBucket struct {
	s3Client   *s3.Client
	bucket     string
	endpoint   string
	folderName string
}

func New() *SpaceBucket {
	bucket := os.Getenv("STORAGE_BUCKET_NAME")
	region := os.Getenv("STORAGE_REGION")
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	folderName := os.Getenv("STORAGE_FOLDER_NAME")

	if strings.TrimSpace(bucket) == "" ||
		strings.TrimSpace(region) == "" ||
		strings.TrimSpace(endpoint) == "" ||
		strings.TrimSpace(folderName) == "" {
		log.Fatalln("Space bucket invalid configuration")
	}

	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	secretKey := os.Getenv("STORAGE_SECRET_KEY")
	if accessKey == "" || secretKey == "" {
		log.Fatalf("DigitalOcean Spaces credentials not set in environment variables")
	}

	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		})),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	return &SpaceBucket{
		s3Client:   s3Client,
		bucket:     bucket,
		endpoint:   endpoint,
		folderName: folderName,
	}
}

func (s *SpaceBucket) UploadImage(ctx echo.Context, image *entity.Image, userID uuid.UUID) (string, *echo.HTTPError) {
	b64data := image.File[strings.IndexByte(image.File, ',')+1:]
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

	image.FileName = uuid.New().String() + fileExtension
	path := fmt.Sprintf("%s/%s/%s/%s", s.folderName, userID.String(), image.Group, image.FileName)

	// Upload the file
	_, errl := s.uploadFile(path, decodedImage)
	if errl != nil {
		return "", errl
	}

	return path, nil
}

func (s *SpaceBucket) uploadFile(filePath string, decodedImage []byte) (*s3.PutObjectOutput, *echo.HTTPError) {
	log.Println(s.bucket)
	log.Println(filePath)
	output, err := s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
		Body:   bytes.NewReader(decodedImage),
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to upload file %q, %v", filePath, err))
	}

	return output, nil
}
