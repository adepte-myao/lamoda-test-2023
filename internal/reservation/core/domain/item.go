package domain

type ItemID string

type Item struct {
	ID              ItemID  `json:"id"`
	Name            string  `json:"name,omitempty"`
	Size            Size    `json:"size,omitempty"`
	WeightKilograms float64 `json:"weightKilograms,omitempty"`
}

func (item *Item) VolumeM2() float64 {
	return item.Size.LengthMeters * item.Size.WidthMeters * item.Size.HeightMeters
}

type Size struct {
	LengthMeters float64 `json:"lengthMeters"`
	WidthMeters  float64 `json:"widthMeters"`
	HeightMeters float64 `json:"heightMeters"`
}
