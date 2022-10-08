package repository

import (
	"context"
	"database/sql"
	"food-track-be/model"
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

func (r *MealRepository) FindAll() ([]*model.Meal, error) {
	var meals []*model.Meal
	err := r.db.NewSelect().Model(&meals).Scan(r.ctx)
	return meals, err
}

func (r *MealRepository) FindById(id uuid.UUID) (*model.Meal, error) {
	var meal model.Meal
	err := r.db.NewSelect().Model(&meal).Where("id = ?", id).Scan(r.ctx)
	return &meal, err
}

func (r *MealRepository) Create(meal *model.Meal) (sql.Result, error) {
	return r.db.NewInsert().Model(meal).Exec(r.ctx)
}

func (r *MealRepository) Update(meal *model.Meal) (sql.Result, error) {
	return r.db.NewUpdate().Model(meal).Where("id = ?", meal.ID).Exec(r.ctx)
}

func (r *MealRepository) Delete(meal *model.Meal) (sql.Result, error) {
	return r.db.NewDelete().Model(meal).Where("id = ?", meal.ID).Exec(r.ctx)
}

func (r *MealRepository) GetAverageKcalEatenInDateRange(startRange time.Time, endRange time.Time) (float64, error) {
	var result float64
	queryStr := "SELECT AVG(COALESCE(kcal, 0)) FROM food_consumption WHERE meal_id IN (SELECT id FROM meal WHERE date BETWEEN ? AND ?) GROUP BY meal_id"
	queryResult, err := r.db.Query(queryStr, startRange, endRange)
	if err != nil {
		return 0, err
	}
	queryResult.Next()
	err = queryResult.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *MealRepository) GetAverageFoodCostInDateRange(startRange time.Time, endRange time.Time) (float64, error) {
	var result float64
	queryStr := "SELECT AVG(COALESCE(cost, 0)) FROM food_consumption WHERE meal_id IN (SELECT id FROM meal WHERE date BETWEEN ? AND ?) GROUP BY meal_id"
	queryResult, err := r.db.Query(queryStr, startRange, endRange)
	if err != nil {
		return 0, err
	}
	queryResult.Next()
	err = queryResult.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *MealRepository) GetSumFoodCostInDateRange(startRange time.Time, endRange time.Time) (float64, error) {
	var result float64
	queryStr := "SELECT SUM(COALESCE(cost, 0)) FROM food_consumption WHERE meal_id IN (SELECT id FROM meal WHERE date BETWEEN ? AND ?) GROUP BY meal_id"
	queryResult, err := r.db.Query(queryStr, startRange, endRange)
	if err != nil {
		return 0, err
	}
	queryResult.Next()
	err = queryResult.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}
