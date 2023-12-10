package ports

import (
	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/domain"
)

type StorehouseRepository interface {
	GetByID(id domain.StoreHouseID) (domain.StoreHouse, error)
	GetUnreservedItemsByID(id domain.StoreHouseID) (map[domain.ItemID]domain.ItemData, error)
	GetAllAsMap() (map[domain.StoreHouseID]domain.StoreHouse, error)
	UpdateAll(storehouses map[domain.StoreHouseID]domain.StoreHouse) error
}

type ItemsRepository interface {
	GetAllAsMap() (map[domain.ItemID]domain.Item, error)
}

type ReservationRepository interface {
	GetByID(id string) (domain.Reservation, error)
	Save(reservation domain.Reservation) error
	Update(reservation domain.Reservation) error
	Delete(id string) error
}
