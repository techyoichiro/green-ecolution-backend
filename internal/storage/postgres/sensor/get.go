package sensor

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/pkg/errors"
)

func (r *SensorRepository) GetAll(ctx context.Context) ([]*entities.Sensor, error) {
	rows, err := r.store.GetAllSensors(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *SensorRepository) GetByID(ctx context.Context, id int32) (*entities.Sensor, error) {
	row, err := r.store.GetSensorByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}

func (r *SensorRepository) GetStatusByID(ctx context.Context, id int32) (*entities.SensorStatus, error) {
	sensor, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &sensor.Status, nil
}

func (r *SensorRepository) GetSensorByStatus(ctx context.Context, status *entities.SensorStatus) ([]*entities.Sensor, error) {
	row, err := r.store.GetSensorByStatus(ctx, sqlc.SensorStatus(*status))
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlList(row), nil
}

func (r *SensorRepository) GetSensorDataByID(ctx context.Context, id int32) ([]*entities.SensorData, error) {
	rows, err := r.store.GetSensorDataBySensorID(ctx, id)
	if err != nil {
		return nil, err
	}

	domainData := make([]*entities.SensorData, len(rows))

	for i, row := range rows {
		domainData[i] = r.mapper.FromSqlSensorData(row)
		data, err := mapper.MapSensorData(row.Data)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map sensor data")
		}
		domainData[i].Data = data
	}

	return domainData, nil
}
