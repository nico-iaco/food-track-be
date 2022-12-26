package dto

type MealStatisticsDto struct {
	AverageWeekCalories            float64                 `json:"averageWeekCalories"`
	AverageWeekCaloriesPerMealType []AvgKcalPerMealTypeDto `json:"averageWeekCaloriesPerMealType"`
	AverageWeekFoodCost            float64                 `json:"averageWeekFoodCost"`
	SumWeekFoodCost                float64                 `json:"sumWeekFoodCost"`
	MostConsumedFood               MostConsumedFoodDto     `json:"mostConsumedFood"`
}
