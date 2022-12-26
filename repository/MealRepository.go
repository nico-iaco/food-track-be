package repository

import (
	"context"
	"database/sql"
	"food-track-be/model"
	"food-track-be/model/dto"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type MealRepository struct {
	db  bun.DB
	ctx context.Context
}

func NewMealRepository(db bun.DB) *MealRepository {
	return &MealRepository{db: db, ctx: context.Background()}
}

func (r *MealRepository) FindAll(userId string) ([]*model.Meal, error) {
	var meals []*model.Meal
	err := r.db.NewSelect().Model(&meals).Where("user_id = ?", userId).Scan(r.ctx)
	return meals, err
}

func (r *MealRepository) FindByIdAndUserId(id uuid.UUID, userId string) (*model.Meal, error) {
	var meal model.Meal
	err := r.db.NewSelect().Model(&meal).Where("id = ?", id).Where("user_id = ?", userId).Scan(r.ctx)
	return &meal, err
}

func (r *MealRepository) Create(meal *model.Meal) (sql.Result, error) {
	return r.db.NewInsert().Model(meal).Exec(r.ctx)
}

func (r *MealRepository) Update(meal *model.Meal, userId string) (sql.Result, error) {
	return r.db.NewUpdate().Model(meal).Where("id = ?", meal.ID).Where("user_id = ?", userId).Exec(r.ctx)
}

func (r *MealRepository) Delete(meal *model.Meal, userId string) (sql.Result, error) {
	return r.db.NewDelete().Model(meal).Where("id = ? ", meal.ID, userId).Where("user_id = ?", userId).Exec(r.ctx)
}

func (r *MealRepository) GetAverageKcalEatenInDateRange(startRange time.Time, endRange time.Time, userId string) (float64, error) {
	var result float64

	endRange = setEndOfTheDay(endRange)

	queryStr := "SELECT COALESCE(SUM(kcal), 0.0) FROM food_consumption WHERE meal_id IN (SELECT id FROM meal WHERE user_id = ? AND date BETWEEN ? AND ?)"
	queryResult, err := r.db.Query(queryStr, userId, startRange, endRange)
	defer queryResult.Close()
	if err != nil {
		return 0, err
	}
	queryResult.Next()
	err = queryResult.Scan(&result)
	queryResult.Close()
	if err != nil {
		return 0, err
	}
	if startRange.Equal(endRange) {
		return result, nil
	}
	rangeInDays := endRange.Sub(startRange).Hours() / 24
	return result / rangeInDays, nil
}

func (r *MealRepository) GetAverageKcalEatenInDateRangePerMealType(startRange time.Time, endRange time.Time, userId string) ([]dto.AvgKcalPerMealTypeDto, error) {
	var result = make([]dto.AvgKcalPerMealTypeDto, 0)
	var rangeInDays float64
	if startRange.Equal(endRange) {
		rangeInDays = 1
	} else {
		rangeInDays = endRange.Sub(startRange).Hours() / 24
	}

	endRange = setEndOfTheDay(endRange)

	queryStr := "SELECT m.meal_type, COALESCE(SUM(kcal), 0.0) / ? as avg_kcal FROM meal m join food_consumption fc on m.id = fc.meal_id WHERE m.user_id = ? AND date BETWEEN ? AND ? group by m.meal_type"
	queryResult, err := r.db.Query(queryStr, rangeInDays, userId, startRange, endRange)

	if err != nil {
		return []dto.AvgKcalPerMealTypeDto{}, err
	}
	for queryResult.Next() {
		var e dto.AvgKcalPerMealTypeDto
		err = queryResult.Scan(&e.MealType, &e.AvgKcal)
		if err != nil {
			return []dto.AvgKcalPerMealTypeDto{}, err
		}
		result = append(result, e)
	}
	queryResult.Close()
	if err != nil {
		return []dto.AvgKcalPerMealTypeDto{}, err
	}
	return result, nil
}

func (r *MealRepository) GetAverageFoodCostInDateRange(startRange time.Time, endRange time.Time, userId string) (float64, error) {
	var result float64

	endRange = setEndOfTheDay(endRange)

	queryStr := "SELECT COALESCE(SUM(cost), 0.0) FROM food_consumption WHERE meal_id IN (SELECT id FROM meal WHERE user_id = ? AND date BETWEEN ? AND ?)"
	queryResult, err := r.db.Query(queryStr, userId, startRange, endRange)
	if err != nil {
		return 0, err
	}
	queryResult.Next()
	err = queryResult.Scan(&result)
	queryResult.Close()
	if err != nil {
		return 0, err
	}
	if startRange.Equal(endRange) {
		return result, nil
	}
	rangeInDays := endRange.Sub(startRange).Hours() / 24
	return result / rangeInDays, nil
}

func (r *MealRepository) GetSumFoodCostInDateRange(startRange time.Time, endRange time.Time, userId string) (float64, error) {
	var result float64

	endRange = setEndOfTheDay(endRange)

	queryStr := "SELECT COALESCE(SUM(cost), 0.0) FROM food_consumption WHERE meal_id IN (SELECT id FROM meal WHERE user_id = ? AND date BETWEEN ? AND ?)"
	queryResult, err := r.db.Query(queryStr, userId, startRange, endRange)
	if err != nil {
		return 0, err
	}
	queryResult.Next()
	err = queryResult.Scan(&result)
	queryResult.Close()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *MealRepository) GetMealInDateRange(startRange time.Time, endRange time.Time, userId string) ([]model.Meal, error) {
	var meals []model.Meal

	endRange = setEndOfTheDay(endRange)

	err := r.db.NewSelect().Model(&meals).Where("date BETWEEN ? AND ?", startRange, endRange).Where("user_id = ?", userId).Order("date ASC").Scan(r.ctx)
	if err != nil {
		return []model.Meal{}, err
	}
	return meals, nil
}

func setEndOfTheDay(t time.Time) time.Time {
	y, m, d := t.Date()
	endOfTheDay := time.Date(y, m, d, 23, 59, 59, 0, t.Location())
	return endOfTheDay
}
