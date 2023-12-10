package domain

type StoreHouseID string

func (id StoreHouseID) IsEmpty() bool {
	return string(id) == ""
}

type StoreHouse struct {
	ID        StoreHouseID
	Name      string
	Location  Location
	ItemsData map[ItemID]ItemData
}

type ItemData struct {
	Item  Item `json:"item"`
	Count int  `json:"count"`
}
