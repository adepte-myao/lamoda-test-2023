package ports

import (
	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/domain"
)

type ReservationService interface {
	Reserve(request domain.ReserveRequest) (ReservationResponseDTO, error)
	Release(reservationID string, itemsToRelease []domain.ReserveEntry) (ReservationResponseDTO, error)
	GetUnreserved(storehouseID domain.StoreHouseID) ([]domain.ItemData, error)
}
