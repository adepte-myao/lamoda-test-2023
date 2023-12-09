package domain

import (
	"math"
)

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

// Haversine formula is used
func getDistance(location1, location2 Location) float64 {
	deltaLat := deg2rad(location2.Latitude - location1.Latitude)
	deltaLon := deg2rad(location2.Longitude - location1.Longitude)

	lat1 := deg2rad(location1.Latitude)
	lat2 := deg2rad(location2.Latitude)

	a := hav(deltaLat) + math.Cos(lat1)*math.Cos(lat2)*hav(deltaLon)
	return 2 * earthR * math.Asin(math.Sqrt(a))
}

const (
	earthR float64 = 6371
)

func hav(x float64) float64 {
	return math.Pow(math.Sin(x/2), 2)
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180
}
