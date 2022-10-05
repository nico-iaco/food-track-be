package dto

import "github.com/google/uuid"

type FoodConsumptionDto struct {
	ID              uuid.UUID `json:"id"`
	MealID          uuid.UUID `json:"mealId"`
	FoodId          uuid.UUID `json:"foodId"`
	TransactionId   uuid.UUID `json:"transactionId"`
	FoodName        string    `json:"foodName"`
	QuantityUsed    float32   `json:"quantityUsed"`
	QuantityUsedStd float32   `json:"quantityUsedStd"`
	Unit            string    `json:"unit"`
	Kcal            float32   `json:"kcal"`
}
