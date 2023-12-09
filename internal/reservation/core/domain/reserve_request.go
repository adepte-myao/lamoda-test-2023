package domain

type ReserveRequest struct {
	DestinationLocation Location       `json:"destinationLocation"`
	ItemsToReserve      []ReserveEntry `json:"itemsToReserve"`
}
