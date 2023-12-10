package services

import (
	"fmt"

	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/domain"
	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/ports"
)

type Service struct {
	storehouseRepo  ports.StorehouseRepository
	itemsRepo       ports.ItemsRepository
	reservationRepo ports.ReservationRepository
}

func (service Service) Reserve(request domain.ReserveRequest) (ports.ReservationResponseDTO, error) {
	storehouses, err := service.storehouseRepo.GetAllAsMap()
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("reserve: receiving storehouses: %w", err)
	}

	reservation, err := domain.NewReservationFromReserveRequest(request, storehouses)
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("reserve: building reservation: %w", err)
	}

	updatedStorehouses, err := reservation.GetUpdatedStorehouses(storehouses, domain.Reserve)
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("reserve: calculating storehouses state: %w", err)
	}

	err = service.storehouseRepo.UpdateAll(updatedStorehouses)
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("reserve: updating storehouses state: %w", err)
	}

	err = service.reservationRepo.Save(reservation)
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("reserve: saving reservation: %w", err)
	}

	items, err := service.itemsRepo.GetAllAsMap()
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("reserve: receiving items: %w", err)
	}

	totalCost, err := reservation.GetTotalCost(storehouses, items)
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("reserve: calculating total cost: %w", err)
	}

	return ports.ReservationResponseDTO{Reservation: reservation, TotalCost: totalCost}, nil
}

func (service Service) Release(reservationID string, itemsToRelease []domain.ReserveEntry) (ports.ReservationResponseDTO, error) {
	reservation, err := service.reservationRepo.GetByID(reservationID)
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("release: receiving reservation: %w", err)
	}

	storehouses, err := service.storehouseRepo.GetAllAsMap()
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("release: receiving storehouses: %w", err)
	}

	oldStorehousesState, err := reservation.GetUpdatedStorehouses(storehouses, domain.Release)
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("release: calculating released storehouse state: %w", err)
	}

	if len(itemsToRelease) > 0 {
		err = reservation.Release(itemsToRelease)
		if err != nil {
			return ports.ReservationResponseDTO{}, fmt.Errorf("release: calculating new reservation state: %w", err)
		}
	}

	needToDeleteReservation := len(itemsToRelease) == 0 || len(reservation.Entries) == 0

	if needToDeleteReservation {
		err = service.reservationRepo.Delete(reservationID)
		if err != nil {
			return ports.ReservationResponseDTO{}, fmt.Errorf("release: deleting reservation: %w", err)
		}
	} else {
		err = service.reservationRepo.Update(reservation)
		if err != nil {
			return ports.ReservationResponseDTO{}, fmt.Errorf("relese: updating reservation: %w", err)
		}
	}

	var newStorehousesState map[domain.StoreHouseID]domain.StoreHouse
	if needToDeleteReservation {
		newStorehousesState = oldStorehousesState
	} else {
		newStorehousesState, err = reservation.GetUpdatedStorehouses(oldStorehousesState, domain.Reserve)
		if err != nil {
			return ports.ReservationResponseDTO{}, fmt.Errorf("release: calculating new storehouses state: %w", err)
		}
	}

	err = service.storehouseRepo.UpdateAll(newStorehousesState)
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("release: updating storehouse state: %w", err)
	}

	items, err := service.itemsRepo.GetAllAsMap()
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("release: receiving items: %w", err)
	}

	totalCost, err := reservation.GetTotalCost(storehouses, items)
	if err != nil {
		return ports.ReservationResponseDTO{}, fmt.Errorf("release: calculating total cost: %w", err)
	}

	return ports.ReservationResponseDTO{Reservation: reservation, TotalCost: totalCost}, nil
}

func (service Service) GetUnreserved(storehouseID domain.StoreHouseID, itemIDs []domain.ItemID) ([]domain.ItemData, error) {
	allUnreserved, err := service.storehouseRepo.GetUnreservedItemsByID(storehouseID)
	if err != nil {
		return nil, fmt.Errorf("get unreserved: receiving all unreserved: %w", err)
	}

	unreserved := make([]domain.ItemData, 0)
	for _, itemID := range itemIDs {
		itemData, ok := allUnreserved[itemID]
		if !ok {
			// TODO: should get any info about missing ones?
			continue
		}

		unreserved = append(unreserved, itemData)
	}

	return unreserved, nil
}
