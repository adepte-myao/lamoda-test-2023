package domain

type ReserveEntry struct {
	ItemID             ItemID       `json:"itemID" validate:"required"`
	Count              int          `json:"count" validate:"required"`
	SourceStorehouseID StoreHouseID `json:"sourceStorehouseID"`
}
