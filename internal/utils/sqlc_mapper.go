package utils

import (
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func PgTimestampToTime(t pgtype.Timestamp) time.Time {
	return t.Time
}

//nolint:gocritic
func ConvertNullableImage(img sqlc.Image) *entities.Image {
	if img.ID == 0 {
		return nil
	}

	return &entities.Image{
		ID:        img.ID,
		CreatedAt: img.CreatedAt.Time,
		UpdatedAt: img.UpdatedAt.Time,
		URL:       img.Url,
	}
}