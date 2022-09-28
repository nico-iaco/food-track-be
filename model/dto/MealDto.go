package dto

import "time"

type MealDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MealType    string `json:"mealType"`
	//FoodTypes   []string  `json:"foodTypes"`
	Date time.Time `json:"date"`
}
