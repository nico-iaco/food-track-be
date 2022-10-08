package dto

type AvgKcalPerMealTypeDto struct {
	MealType string  `json:"mealType"`
	AvgKcal  float64 `json:"avgKcal"`
}
