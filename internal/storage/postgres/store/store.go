package store

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type EntityType string

const (
	Sensor      EntityType = "sensor"
	Image       EntityType = "image"
	Flowerbed   EntityType = "flowerbed"
	TreeCluster EntityType = "treecluster"
	Tree        EntityType = "tree"
	Vehicle     EntityType = "vehicle"
)

type Store struct {
	*sqlc.Queries
	db         *pgx.Conn
	entityType EntityType
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		Queries: sqlc.New(db),
		db:      db,
	}
}

func (s *Store) SetEntityType(entityType EntityType) {
	s.entityType = entityType
}

func (s *Store) HandleError(err error) error {
	if err == nil {
		return nil
	}

	slog.Error("An Error occurred in database operation", "error", err, "entityType", s.entityType)
	switch err {
	case pgx.ErrNoRows:
		switch s.entityType {
		case Image:
			slog.Error("Image not found", "error", err, "stack", errors.WithStack(err))
			return storage.ErrImageNotFound
		case Sensor:
			slog.Error("Sensor not found", "error", err, "stack", errors.WithStack(err))
			return storage.ErrSensorNotFound
		case Flowerbed:
			slog.Error("Flowerbed not found", "error", err, "stack", errors.WithStack(err))
			return storage.ErrFlowerbedNotFound
		case TreeCluster:
			slog.Error("TreeCluster not found", "error", err, "stack", errors.WithStack(err))
			return storage.ErrTreeClusterNotFound
		default:
			slog.Error("Entity not found", "error", err, "stack", errors.WithStack(err))
			return storage.ErrEntityNotFound
		}
	case pgx.ErrTooManyRows:
		slog.Error("Receive more rows then expected", "error", err, "stack", errors.WithStack(err))
		return storage.ErrToManyRows
	case pgx.ErrTxClosed:
		slog.Error("Connection is closed", "error", err, "stack", errors.WithStack(err))
		return storage.ErrTxClosed
	case pgx.ErrTxCommitRollback:
		slog.Error("Transaction cannot commit or rollback", "error", err, "stack", errors.WithStack(err))
		return storage.ErrTxCommitRollback
	case sql.ErrConnDone:
		slog.Error("Connection is closed", "error", err, "stack", errors.WithStack(err))
		return storage.ErrConnectionClosed

	default:
		slog.Error("Unknown error", "error", err, "stack", errors.WithStack(err))
		return errors.Wrap(err, "unknown error in postgres store")
	}
}

func (s *Store) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *Store) Close() {
	s.db.Close(context.Background())
}

func (s *Store) CheckSensorExists(ctx context.Context, sensorID *int32) error {
	if sensorID != nil {
		_, err := s.GetSensorByID(ctx, *sensorID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return storage.ErrSensorNotFound
			} else {
				slog.Error("Error getting sensor by id", "error", err)
				return s.HandleError(err)
			}
		}
	}

	return nil
}

func (s *Store) CheckImageExists(ctx context.Context, imageID *int32) error {
	if imageID != nil {
		_, err := s.GetImageByID(ctx, *imageID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return storage.ErrImageNotFound
			} else {
				slog.Error("Error getting image by id", "error", err)
				return s.HandleError(err)
			}
		}
	}

	return nil
}
