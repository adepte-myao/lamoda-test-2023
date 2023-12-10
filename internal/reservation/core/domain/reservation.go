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
	ErrUnknownItem                    = errors.New("unknown item")
	ErrNotEnoughItemsInStorehouse     = errors.New("not enough items in storehouse")
	ErrNotEnoughItemsInAllStorehouses = errors.New("not enough items in all storehouses")
	ErrInvalidReleaseItems            = errors.New("invalid release items")
	ErrNotEnoughItemsInReservation    = errors.New("not enough items in reservation")
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
func (reservation *Reservation) GetTotalCost(storehouses map[StoreHouseID]StoreHouse, items map[ItemID]Item) (float64, error) {
	storehouseEntries := groupEntriesPerStorehouse(reservation.Entries)

	var totalCost float64 = 0
	var resultErr error
	for storehouseID, entries := range storehouseEntries {
		storehouse, ok := storehouses[storehouseID]
		if !ok {
			resultErr = errors.Join(resultErr, fmt.Errorf("%w: %s", ErrUnknownStorehouse, storehouseID))
			continue
		}

		distance := getDistance(storehouse.Location, reservation.DestinationLocation)

		var totalEntryMetric float64 = 0
		for _, entry := range entries {
			item, ok := items[entry.ItemID]
			if !ok {
				resultErr = errors.Join(resultErr, fmt.Errorf("%w: %s", ErrUnknownItem, entry.ItemID))
				continue
			}

			totalEntryMetric += math.Log(max(item.WeightKilograms, item.VolumeM2())) * float64(entry.Count)
		}

		totalCost += k1*distance*totalEntryMetric + k2
	}

	return totalCost, resultErr
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

func (reservation *Reservation) GetUpdatedStorehouses(oldStorehouses map[StoreHouseID]StoreHouse, op OperationType, items map[ItemID]Item) (map[StoreHouseID]StoreHouse, error) {
	updatedSH := maps.Clone(oldStorehouses)
	for storehouseID, storehouse := range oldStorehouses {
		updatedStorehouse := updatedSH[storehouseID]
		updatedStorehouse.ItemsData = maps.Clone(storehouse.ItemsData)
		updatedSH[storehouseID] = updatedStorehouse
	}

	for _, entry := range reservation.Entries {
		storehouse, ok := updatedSH[entry.SourceStorehouseID]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownStorehouse, entry.SourceStorehouseID)
		}

		itemData, ok := storehouse.ItemsData[entry.ItemID]
		if !ok && op == Reserve {
			return nil, fmt.Errorf("%w: %s", ErrUnknownItem, entry.ItemID)
		}
		// TODO: make simpler
		if itemInfo, itemExists := items[entry.ItemID]; op == Release && !ok && itemExists {
			itemData = ItemData{
				Item:  itemInfo,
				Count: 0,
			}
		}

		if op == Reserve {
			if itemData.Count < entry.Count {
				return nil, fmt.Errorf("%w: expected at least: %d, have: %d", ErrNotEnoughItemsInStorehouse, entry.Count, itemData.Count)
			}

			itemData.Count -= entry.Count
		} else if op == Release {
			itemData.Count += entry.Count
		}

		if itemData.Count == 0 {
			delete(storehouse.ItemsData, entry.ItemID)
		} else {
			storehouse.ItemsData[entry.ItemID] = itemData
		}

		updatedSH[entry.SourceStorehouseID] = storehouse
	}

	return updatedSH, nil
}

func NewReservationFromReserveRequest(request ReserveRequest, storehouses map[StoreHouseID]StoreHouse) (Reservation, error) {
	reservation := Reservation{
		ID:                  uuid.New().String(),
		DestinationLocation: request.DestinationLocation,
		Entries:             make([]ReserveEntry, 0),
	}

	var resultErr error

	knownEntries, leftEntries, updatedStorehouses, err := filterKnownDistributions(request.ItemsToReserve, storehouses)
	resultErr = errors.Join(resultErr, err)

	reservation.Entries = knownEntries

	sortedStorehouses := sortStorehousesByDistance(updatedStorehouses, request.DestinationLocation)

	distributedEntries, err := distribute(leftEntries, sortedStorehouses)
	resultErr = errors.Join(resultErr, err)

	reservation.Entries = append(reservation.Entries, distributedEntries...)

	return reservation, resultErr
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

func distribute(entries []ReserveEntry, sortedStorehouses []StoreHouse) (distributed []ReserveEntry, err error) {
	for _, entry := range entries {
		// using greedy algorithm for each entry: just take all required items from the nearest left storehouse
		// till either it's enough items or no storehouses left

		for i, storehouse := range sortedStorehouses {
			if itemData, ok := storehouse.ItemsData[entry.ItemID]; ok {
				if itemData.Count >= entry.Count {
					itemData.Count -= entry.Count
					sortedStorehouses[i].ItemsData[entry.ItemID] = itemData

					distributed = append(distributed, ReserveEntry{
						ItemID:             entry.ItemID,
						Count:              entry.Count,
						SourceStorehouseID: storehouse.ID,
					})

					entry.Count = 0
				} else {
					entry.Count -= itemData.Count

					distributed = append(distributed, ReserveEntry{
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
			err = errors.Join(err,
				fmt.Errorf("%w, item: %s", ErrNotEnoughItemsInAllStorehouses, entry.ItemID))
		}
	}

	return distributed, err
}

func (reservation *Reservation) Release(items []ReserveEntry) error {
	// TODO: handle cases when storehouse id is not provided

	type key struct {
		ItemID       ItemID
		StorehouseID StoreHouseID
	}

	toRelease := make(map[key]ReserveEntry)
	for _, item := range items {
		toRelease[key{ItemID: item.ItemID, StorehouseID: item.SourceStorehouseID}] = item
	}

	for i, entry := range reservation.Entries {
		itemKey := key{ItemID: entry.ItemID, StorehouseID: entry.SourceStorehouseID}
		itemToRelease, ok := toRelease[itemKey]
		if !ok {
			// this item is not in release list
			continue
		}

		if entry.Count < itemToRelease.Count {
			return fmt.Errorf("%w: expected at least %d, got %d", ErrNotEnoughItemsInReservation, itemToRelease.Count, entry.Count)
		}

		reservation.Entries[i].Count -= itemToRelease.Count
		delete(toRelease, itemKey)
	}

	if len(toRelease) > 0 {
		return fmt.Errorf("%w, elements: %v", ErrInvalidReleaseItems, toRelease)
	}

	// get rid of all reservation items that has count == 0
	i := 0
	for i < len(reservation.Entries) {
		if reservation.Entries[i].Count == 0 {
			reservation.Entries = append(reservation.Entries[:i], reservation.Entries[:i+1]...)
		} else {
			i++
		}
	}

	return nil
}
