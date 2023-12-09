package domain

type ItemID string

type Item struct {
	ID              ItemID
	Name            string
	Size            Size
	WeightKilograms float64
}

func (item *Item) VolumeM2() float64 {
	return item.Size.LengthMeters * item.Size.WidthMeters * item.Size.HeightMeters
}

type Size struct {
	LengthMeters float64
	WidthMeters  float64
	HeightMeters float64
}
