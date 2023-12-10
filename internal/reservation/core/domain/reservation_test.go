package domain

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReservation_GetTotalCost(t *testing.T) {
	_, storehouses := getReservedItemsAndStorehouses()
	items := getItems()

	reservation := Reservation{
		ID:                  "",
		DestinationLocation: Location{50, 50},
		Entries:             getExpectedEntries(),
	}

	reservation.Entries = append(reservation.Entries,
		ReserveEntry{ItemID: "10", Count: 3, SourceStorehouseID: "a"},  // unknown item
		ReserveEntry{ItemID: "2", Count: 0, SourceStorehouseID: "AAA"}, // unknown storehouse
	)

	cost, err := reservation.GetTotalCost(storehouses, items)
	if !assert.Error(t, err) {
		t.FailNow()
	}

	expectedErrors := []error{ErrUnknownStorehouse, ErrUnknownItem}
	for _, expectedErr := range expectedErrors {
		assert.True(t, errors.Is(err, expectedErr))
	}

	var expectedCost float64 = 65500
	assert.Less(t, math.Abs(expectedCost-cost), 100.0)
}

func TestNewReservationFromReserveRequest(t *testing.T) {
	itemsToReserve, storehouses := getReservedItemsAndStorehouses()

	request := ReserveRequest{
		DestinationLocation: Location{40, 40},
		ItemsToReserve:      itemsToReserve,
	}

	reservation, err := NewReservationFromReserveRequest(request, storehouses)

	expectedErr := errors.Join(
		fmt.Errorf("%w: storehouse id: %s, item id: %d", ErrNotEnoughItemsInStorehouse, "a", 1),
		fmt.Errorf("%w: storehouse id: %s, item id: %d", ErrNotEnoughItemsInStorehouse, "a", 2),
		fmt.Errorf("%w: %s", ErrUnknownStorehouse, "A"),
		fmt.Errorf("%w, item: %d", ErrNotEnoughItemsInAllStorehouses, 4),
		fmt.Errorf("%w, item: %d", ErrNotEnoughItemsInAllStorehouses, 5),
	)
	assert.EqualValues(t, expectedErr.Error(), err.Error())

	expectedEntries := getExpectedEntries()
	assert.EqualValues(t, expectedEntries, reservation.Entries)
}

func getReservedItemsAndStorehouses() ([]ReserveEntry, map[StoreHouseID]StoreHouse) {
	items := []ReserveEntry{
		{ItemID: "1", Count: 5, SourceStorehouseID: "a"}, // Storehouse filled, it's not enough items
		{ItemID: "2", Count: 5, SourceStorehouseID: "a"}, // Storehouse filled, no items of the type there
		{ItemID: "3", Count: 5, SourceStorehouseID: "a"}, // Storehouse filled, it's enough items
		{ItemID: "4", Count: 5, SourceStorehouseID: ""},  // Storehouse not filled, no items of the type at all
		{ItemID: "5", Count: 5, SourceStorehouseID: ""},  // Storehouse not filled, not enough items on all storehouses
		{ItemID: "6", Count: 5, SourceStorehouseID: ""},  // Storehouse not filled, enough items on the nearest storehouse
		{ItemID: "7", Count: 5, SourceStorehouseID: ""},  // Storehouse not filled, enough items on non-nearest storehouse
		{ItemID: "8", Count: 5, SourceStorehouseID: ""},  // Storehouse not filled, enough items on several storehouse (in total)
		{ItemID: "9", Count: 5, SourceStorehouseID: "A"}, // Storehouse filled, but it does not exist
	}

	sizes := make([]Size, len(items))
	for i := range sizes {
		sizes[i] = getRandomSize()
	}

	storehouses := map[StoreHouseID]StoreHouse{
		"a": {ID: "a", Name: "a", Location: Location{Latitude: 50, Longitude: 50}, ItemsData: map[ItemID]ItemData{
			"1": {Item: Item{ID: "1", Name: "1", Size: sizes[0], WeightKilograms: 1}, Count: 2},
			"3": {Item: Item{ID: "3", Name: "3", Size: sizes[2], WeightKilograms: 3}, Count: 7},
			"5": {Item: Item{ID: "5", Name: "5", Size: sizes[4], WeightKilograms: 5}, Count: 1},
			"6": {Item: Item{ID: "6", Name: "6", Size: sizes[5], WeightKilograms: 6}, Count: 6},
			"8": {Item: Item{ID: "8", Name: "8", Size: sizes[7], WeightKilograms: 8}, Count: 3},
		}},
		"b": {ID: "b", Name: "b", Location: Location{Latitude: 60, Longitude: 60}, ItemsData: map[ItemID]ItemData{
			"5": {Item: Item{ID: "5", Name: "5", Size: sizes[4], WeightKilograms: 5}, Count: 1},
			"7": {Item: Item{ID: "7", Name: "7", Size: sizes[6], WeightKilograms: 7}, Count: 6},
			"8": {Item: Item{ID: "8", Name: "8", Size: sizes[7], WeightKilograms: 8}, Count: 3},
		}},
	}

	return items, storehouses
}

func getItems() map[ItemID]Item {
	sizes := make([]Size, 8)
	for i := range sizes {
		sizes[i] = Size{
			LengthMeters: float64(i + 1),
			WidthMeters:  float64(i + 2),
			HeightMeters: float64(i + 3),
		}
	}

	items := make(map[ItemID]Item, 8)
	for i := 0; i < len(sizes); i++ {
		itemID := ItemID(strconv.Itoa(i + 1))
		items[itemID] = Item{
			ID:              itemID,
			Name:            string(itemID),
			Size:            sizes[i],
			WeightKilograms: float64(i + 1),
		}
	}

	return items
}

func getExpectedEntries() []ReserveEntry {
	return []ReserveEntry{
		{ItemID: "3", Count: 5, SourceStorehouseID: "a"},
		{ItemID: "5", Count: 1, SourceStorehouseID: "a"},
		{ItemID: "5", Count: 1, SourceStorehouseID: "b"},
		{ItemID: "6", Count: 5, SourceStorehouseID: "a"},
		{ItemID: "7", Count: 5, SourceStorehouseID: "b"},
		{ItemID: "8", Count: 3, SourceStorehouseID: "a"},
		{ItemID: "8", Count: 2, SourceStorehouseID: "b"},
	}
}

func getRandomSize() Size {
	return Size{
		LengthMeters: rand.Float64() * 10,
		WidthMeters:  rand.Float64() * 3,
		HeightMeters: rand.Float64() * 2,
	}
}
