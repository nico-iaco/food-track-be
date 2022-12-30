package dto

import "github.com/google/uuid"

type MostConsumedFoodDto struct {
	FoodId          uuid.UUID 		`json:"foodId"`
	FoodName        string 			`json:"foodName"`
	QuantityUsed    float32 		`json:"quantityUsed"`
	QuantityUsedStd float32 		`json:"quantityUsedStd"`
	Unit            string 			`json:"unit"`
}
