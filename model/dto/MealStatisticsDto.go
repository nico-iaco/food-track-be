package dto

type MealStatisticsDto struct {
	AverageWeekCalories float64 `json:"averageWeekCalories"`
	AverageWeekFoodCost float64 `json:"averageWeekFoodCost"`
	SumWeekFoodCost     float64 `json:"sumWeekFoodCost"`
}
