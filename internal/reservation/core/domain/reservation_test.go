package domain

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterKnownDistributions(t *testing.T) {
	entriesToFilter := []ReserveEntry{
		{ItemID: "1", Count: 5, SourceStorehouseID: "a"},
		{ItemID: "2", Count: 1, SourceStorehouseID: "a"},
		{ItemID: "3", Count: 2, SourceStorehouseID: ""},
	}

	sizes := [3]Size{}
	for i := range sizes {
		sizes[i] = getRandomSize()
	}
	storehouses := map[StoreHouseID]StoreHouse{
		"a": {ID: "a", Name: "A great a", Location: Location{Latitude: 50, Longitude: 50}, ItemsData: map[ItemID]ItemData{
			"1": {Item: Item{ID: "1", Name: "Really beautiful 1", Size: sizes[0], WeightKilograms: 0}, Count: 2},
			"2": {Item: Item{ID: "2", Name: "Incomprehensible 2", Size: sizes[1], WeightKilograms: 0}, Count: 2},
			"3": {Item: Item{ID: "3", Name: "Strange 3", Size: sizes[2], WeightKilograms: 0}, Count: 2},
		}},
	}

	known, left, updatedStorehouses, err := filterKnownDistributions(entriesToFilter, storehouses)

	assert.Error(t, err)

	knownShouldBe := []ReserveEntry{
		{ItemID: "2", Count: 1, SourceStorehouseID: "a"},
	}

	leftShouldBe := []ReserveEntry{
		{ItemID: "3", Count: 2, SourceStorehouseID: ""},
	}

	updatedStorehousesShouldBe := map[StoreHouseID]StoreHouse{
		"a": {ID: "a", Name: "A great a", Location: Location{Latitude: 50, Longitude: 50}, ItemsData: map[ItemID]ItemData{
			"1": {Item: Item{ID: "1", Name: "Really beautiful 1", Size: sizes[0], WeightKilograms: 0}, Count: 2},
			"2": {Item: Item{ID: "2", Name: "Incomprehensible 2", Size: sizes[1], WeightKilograms: 0}, Count: 1},
			"3": {Item: Item{ID: "3", Name: "Strange 3", Size: sizes[2], WeightKilograms: 0}, Count: 2},
		}},
	}

	assert.EqualValues(t, knownShouldBe, known)
	assert.EqualValues(t, leftShouldBe, left)
	assert.EqualValues(t, updatedStorehousesShouldBe, updatedStorehouses)
}

func getRandomSize() Size {
	return Size{
		LengthMeters: rand.Float64() * 10,
		WidthMeters:  rand.Float64() * 3,
		HeightMeters: rand.Float64() * 2,
	}
}
