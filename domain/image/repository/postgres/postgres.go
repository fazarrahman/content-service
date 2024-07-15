package postgres

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fazarrahman/content-service/domain/image/entity"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(db *gorm.DB) *Postgres {
	return &Postgres{db}
}

func (p *Postgres) Save(ctx echo.Context, image *entity.Image) *echo.HTTPError {
	tx := p.db.Begin()
	tags := []entity.Tag{}
	if err := tx.Where("name IN ?", image.TagStr).Find(&tags).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error when getting existing tags : %s", err.Error()))
	}
	tagMap := make(map[string]*entity.Tag)
	for _, t := range tags {
		tagMap[t.Name] = &t
	}
	notExistTag := []entity.Tag{}
	for _, t := range image.TagStr {
		if tagMap[t] == nil {
			notExistTag = append(notExistTag, entity.Tag{
				Name:      t,
				CreatedAt: time.Now(),
			})
		}
	}

	if err := tx.Clauses(clause.Returning{}).Create(&notExistTag).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error when saving tags : %s", err.Error()))
	}

	image.Tags = tags
	image.Tags = append(image.Tags, notExistTag...)
	if err := tx.WithContext(ctx.Request().Context()).Save(image).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error when saving images : %s", err.Error()))
	}
	tx.Commit()
	return nil
}

func (p *Postgres) GetList(ctx echo.Context, page, size int) ([]*entity.Image, *echo.HTTPError) {
	images := []*entity.Image{}
	offset := (page - 1) * size
	tx := p.db.WithContext(ctx.Request().Context()).
		Limit(size).Offset(offset).
		Find(&images)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error when getting image list : ", tx.Error)
	}
	return images, nil
}
