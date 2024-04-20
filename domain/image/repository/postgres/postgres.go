package postgres

import (
	"errors"
	"net/http"

	"github.com/fazarrahman/content-service/domain/image/entity"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(db *gorm.DB) *Postgres {
	return &Postgres{db}
}

func (p *Postgres) Save(ctx echo.Context, image *entity.Image) *echo.HTTPError {
	if err := p.db.WithContext(ctx.Request().Context()).Save(image).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error when saving image : ", err)
	}
	return nil
}

func (p *Postgres) GetList(ctx echo.Context) ([]*entity.Image, *echo.HTTPError) {
	images := []*entity.Image{}
	tx := p.db.WithContext(ctx.Request().Context()).Find(&images)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error when getting image list : ", tx.Error)
	}
	return images, nil
}
