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
	TotalCost   float64            `json:"totalCost"`
}

type GetUnreservedRequestDTO struct {
	StorehouseID domain.StoreHouseID `form:"storehouse-id" validate:"required"`
}

type GetUnreservedResponseDTO struct {
	StorehouseID domain.StoreHouseID `json:"storehouseID"`
	Items        []domain.ItemData   `json:"items"`
}
