package ports

import (
	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/domain"
)

type ReleaseRequestDTO struct {
	ReservationID  string                `json:"reservationID" validate:"required"`
	ItemsToRelease []domain.ReserveEntry `json:"itemsToRelease"`
}

type ReservationResponseDTO struct {
	Reservation domain.Reservation `json:"reservation"`
}

type GetUnreservedRequestDTO struct {
	StorehouseID domain.StoreHouseID `json:"storehouseID" validate:"required"`
	ItemIDs      []domain.ItemID     `json:"itemIDs"`
}

type GetUnreservedResponseDTO struct {
	StorehouseID domain.StoreHouseID `json:"storehouseID"`
	Items        []domain.ItemData   `json:"items"`
}
