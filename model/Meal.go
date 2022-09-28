package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Meal struct {
	bun.BaseModel `bun:"table:meal,alias:m"`
	ID            uuid.UUID `bun:"type:uuid,nullzero,pk"`
	Name          string    `bun:"type:varchar(255),notnull"`
	Description   string    `bun:"type:varchar(255),nullzero"`
	MealType      MealType  `bun:"type:varchar(255),notnull"`
	//FoodTypes        []FoodType         `bun:"type:varchar(255)[]"`
	Date             string             `bun:"type:varchar(255),notnull"`
	FoodConsumptions []*FoodConsumption `bun:"rel:has-many,join:id=meal_id"`
}

type MealType string

const (
	Breakfast MealType = "breakfast"
	Lunch     MealType = "lunch"
	Dinner    MealType = "dinner"
	Others    MealType = "others"
)

type FoodType string

const (
	Sweet    FoodType = "sweet"
	Sour     FoodType = "sour"
	Salty    FoodType = "salty"
	Bitter   FoodType = "bitter"
	Beverage FoodType = "beverage"
)
