package dto

import "github.com/google/uuid"

type FoodAvailableDto struct {
	ID                uuid.UUID `json:"id"`
	Barcode           string    `json:"barcode"`
	Name              string    `json:"name"`
	Quantity          float32   `json:"quantity"`
	AvailableQuantity float32   `json:"availableQuantity"`
	Unit              string    `json:"unit"`
}
