package domain

import (
	"cmp"
	"errors"
	"fmt"
	"maps"
	"math"
	"slices"

	"github.com/google/uuid"
)

var (
	ErrUnknownStorehouse              = errors.New("unknown storehouse")
	ErrNotEnoughItemsInStorehouse     = errors.New("not enough items in storehouse")
	ErrNotEnoughItemsInAllStorehouses = errors.New("not enough items in all storehouses")
)

const (
	k1 float64 = 1
	k2 float64 = 1e3
)

type Reservation struct {
	ID                  string         `json:"id"`
	DestinationLocation Location       `json:"destinationLocation"`
	Entries             []ReserveEntry `json:"entries"`
}

// GetTotalCost returns transport cost of the reservation.
// The transport cost is calculated with the formula:
//
//	k1 * distance_km * ln(max(mass_kg, volume_m2)) + k2
//
// where k1 * ... is added per item and k2 is added per storehouse
func (reservation *Reservation) GetTotalCost(storehouses map[StoreHouseID]StoreHouse, items map[ItemID]Item) float64 {
	// 1. Group items by source storehouse
	storehouseEntries := groupEntriesPerStorehouse(reservation.Entries)

	// 2. Sum all costs per storehouse
	var totalCost float64 = 0
	for storehouse, entries := range storehouseEntries {

		var totalEntryMetric float64 = 0
		for _, entry := range entries {
			// TODO: check if items are not in map
			item := items[entry.ItemID]
			totalEntryMetric += math.Log(max(item.WeightKilograms, item.VolumeM2())) * float64(entry.Count)
		}

		// TODO: check if storehouse is not in map
		distance := getDistance(storehouses[storehouse].Location, reservation.DestinationLocation)

		totalCost += k1*distance*totalEntryMetric + k2
	}

	return totalCost
}

func groupEntriesPerStorehouse(entries []ReserveEntry) map[StoreHouseID][]ReserveEntry {
	storehouseEntries := make(map[StoreHouseID][]ReserveEntry)
	for _, entry := range entries {
		associatedEntries := storehouseEntries[entry.SourceStorehouseID]
		associatedEntries = append(associatedEntries, entry)
		storehouseEntries[entry.SourceStorehouseID] = associatedEntries
	}

	return storehouseEntries
}

func NewReservationFromReserveRequest(request ReserveRequest, storehouses map[StoreHouseID]StoreHouse) (Reservation, error) {
	reservation := Reservation{
		ID:                  uuid.New().String(),
		DestinationLocation: request.DestinationLocation,
		Entries:             make([]ReserveEntry, 0),
	}

	knownEntries, leftEntries, updatedStorehouses, err := filterKnownDistributions(request.ItemsToReserve, storehouses)
	if err != nil {
		return Reservation{}, nil
	}

	reservation.Entries = knownEntries

	sortedStorehouses := sortStorehousesByDistance(updatedStorehouses, request.DestinationLocation)

	var resultErr error
	distributedEntries := make([]ReserveEntry, 0)
	for _, entry := range leftEntries {
		// using greedy algorithm for each entry

		for i, storehouse := range sortedStorehouses {
			if itemData, ok := storehouse.ItemsData[entry.ItemID]; ok {
				if itemData.Count >= entry.Count {
					itemData.Count -= entry.Count
					sortedStorehouses[i].ItemsData[entry.ItemID] = itemData

					distributedEntries = append(distributedEntries, ReserveEntry{
						ItemID:             entry.ItemID,
						Count:              entry.Count,
						SourceStorehouseID: storehouse.ID,
					})

					entry.Count = 0
				} else {
					entry.Count -= itemData.Count

					distributedEntries = append(distributedEntries, ReserveEntry{
						ItemID:             entry.ItemID,
						Count:              itemData.Count,
						SourceStorehouseID: storehouse.ID,
					})

					itemData.Count = 0
					sortedStorehouses[i].ItemsData[entry.ItemID] = itemData
				}
			}

			if entry.Count == 0 {
				break
			}
		}

		if entry.Count > 0 {
			err = fmt.Errorf("%w, item: %s", ErrNotEnoughItemsInAllStorehouses, entry.ItemID)
			resultErr = errors.Join(resultErr, err)
		}
	}

	if resultErr != nil {
		return Reservation{}, resultErr
	}

	reservation.Entries = append(reservation.Entries, distributedEntries...)

	return reservation, nil
}

func filterKnownDistributions(
	entriesToFilter []ReserveEntry, storehouses map[StoreHouseID]StoreHouse) (
	known, left []ReserveEntry, updatedStorehouses map[StoreHouseID]StoreHouse, resultErr error) {

	updatedStorehouses = maps.Clone(storehouses)
	for storehouseID, storehouse := range storehouses {
		storehouseDst := updatedStorehouses[storehouseID]
		storehouseDst.ItemsData = maps.Clone(storehouse.ItemsData)
		updatedStorehouses[storehouseID] = storehouseDst
	}

	for _, entry := range entriesToFilter {
		if entry.SourceStorehouseID.IsEmpty() {
			left = append(left, entry)
			continue
		}

		storehouse, ok := storehouses[entry.SourceStorehouseID]
		if !ok {
			err := fmt.Errorf("%w: %s", ErrUnknownStorehouse, entry.SourceStorehouseID)
			resultErr = errors.Join(resultErr, err)
			continue
		}

		itemData, ok := storehouse.ItemsData[entry.ItemID]
		if !ok || itemData.Count < entry.Count {
			err := fmt.Errorf("%w: storehouse id: %s, item id: %s", ErrNotEnoughItemsInStorehouse, storehouse.ID, entry.ItemID)
			resultErr = errors.Join(resultErr, err)
			continue
		}

		known = append(known, entry)

		itemData.Count -= entry.Count
		storehouse.ItemsData[entry.ItemID] = itemData
		updatedStorehouses[entry.SourceStorehouseID] = storehouse
	}

	return known, left, updatedStorehouses, resultErr
}

func sortStorehousesByDistance(storehouses map[StoreHouseID]StoreHouse, from Location) []StoreHouse {
	slice := make([]StoreHouse, 0, len(storehouses))
	for _, storehouse := range storehouses {
		slice = append(slice, storehouse)
	}

	slices.SortFunc(slice, func(a, b StoreHouse) int {
		return cmp.Compare(getDistance(a.Location, from), getDistance(b.Location, from))
	})

	return slice
}
