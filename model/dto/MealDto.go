package dto

import (
	"food-track-be/model"
	"github.com/google/uuid"
	"time"
)

type MealDto struct {
	ID          uuid.UUID      `json:"id,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	MealType    model.MealType `json:"mealType"`
	Date        time.Time      `json:"date"`
	//FoodTypes   []string  `json:"foodTypes"`
}
