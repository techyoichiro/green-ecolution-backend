package treecluster

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func defaultTreeCluster() *entities.TreeCluster {
	return &entities.TreeCluster{
		Region:         nil,
		Address:        "",
		Description:    "",
		MoistureLevel:  0,
		Latitude:       nil,
		Longitude:      nil,
		WateringStatus: entities.TreeClusterWateringStatusUnknown,
		SoilCondition:  entities.TreeSoilConditionUnknown,
		Archived:       false,
		LastWatered:    nil,
		Trees:          make([]*entities.Tree, 0),
		Name:           "",
	}
}

func (r *TreeClusterRepository) Create(ctx context.Context, tcFn ...entities.EntityFunc[entities.TreeCluster]) (*entities.TreeCluster, error) {
	entity := defaultTreeCluster()
	for _, fn := range tcFn {
		fn(entity)
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}
	entity.ID = id

	return r.GetByID(ctx, id)
}

func (r *TreeClusterRepository) createEntity(ctx context.Context, entity *entities.TreeCluster) (int32, error) {
	var region *int32
	if entity.Region != nil {
		region = &entity.Region.ID
	}

	args := sqlc.CreateTreeClusterParams{
		RegionID:       region,
		Address:        entity.Address,
		Description:    entity.Description,
		MoistureLevel:  entity.MoistureLevel,
		WateringStatus: sqlc.TreeClusterWateringStatus(entity.WateringStatus),
		SoilCondition:  sqlc.TreeSoilCondition(entity.SoilCondition),
		Name:           entity.Name,
	}

	id, err := r.store.CreateTreeCluster(ctx, &args)
	if err != nil {
		return -1, err
	}

	if entity.Latitude != nil && entity.Longitude != nil {
		err = r.store.SetTreeClusterLocation(ctx, &sqlc.SetTreeClusterLocationParams{
			ID:        id,
			Latitude:  entity.Latitude,
			Longitude: entity.Longitude,
		})
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}
