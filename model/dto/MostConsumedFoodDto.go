package dto

import "github.com/google/uuid"

type MostConsumedFoodDto struct {
	FoodId          uuid.UUID
	FoodName        string
	QuantityUsed    float32
	QuantityUsedStd float32
	Unit            string
}
