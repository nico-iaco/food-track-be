package dto

import (
	"github.com/google/uuid"
)

type FoodTransactionDto struct {
	ID                uuid.UUID `json:"id,omitempty"`
	Vendor            string    `json:"vendor"`
	Quantity          float32   `json:"quantity"`
	AvailableQuantity float32   `json:"availableQuantity"`
	Unit              string    `json:"unit"`
	Price             float32   `json:"price"`
	//ExpirationDate    time.Time `json:"expirationDate"`
}
