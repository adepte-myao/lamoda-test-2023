package ports

import (
	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/domain"
)

type StorehouseRepository interface {
	GetStorehouseByID(id domain.StoreHouseID) domain.StoreHouse
	GetAllStorehousesWithoutItemsAsMap() map[domain.StoreHouseID]domain.StoreHouse
}

type ItemsRepository interface {
	GetAllAsMap() map[domain.ItemID]domain.Item
}
