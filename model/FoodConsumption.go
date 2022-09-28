package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type FoodConsumption struct {
	bun.BaseModel   `bun:"table:food_consumption,alias:fc"`
	ID              uuid.UUID `bun:"type:uuid,notnull,pk,default:uuid_generate_v4()"`
	MealID          uuid.UUID
	FoodId          uuid.UUID
	TransactionId   uuid.UUID
	QuantityUsed    float32
	QuantityUsedStd float32
	Unit            string
	Kcal            float32
}
