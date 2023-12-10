package ports

import (
	"context"

	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/domain"
)

type StorehouseRepository interface {
	GetItemsByID(ctx context.Context, id domain.StoreHouseID) (map[domain.ItemID]domain.ItemData, error)
	GetAllAsMap(ctx context.Context) (map[domain.StoreHouseID]domain.StoreHouse, error)
	UpdateAll(ctx context.Context, storehouses map[domain.StoreHouseID]domain.StoreHouse) error
}

type ItemsRepository interface {
	GetAllAsMap(ctx context.Context) (map[domain.ItemID]domain.Item, error)
}

type ReservationRepository interface {
	GetByID(ctx context.Context, id string) (domain.Reservation, error)
	Save(ctx context.Context, reservation domain.Reservation) error
	Update(ctx context.Context, reservation domain.Reservation) error
	Delete(ctx context.Context, id string) error
}
