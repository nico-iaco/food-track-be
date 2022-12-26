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

type FoodConsumptionRepository struct {
	db  bun.DB
	ctx context.Context
}

func NewFoodConsumptionRepository(db bun.DB) *FoodConsumptionRepository {
	return &FoodConsumptionRepository{db: db, ctx: context.Background()}
}

func (r *FoodConsumptionRepository) FindAll() ([]*model.FoodConsumption, error) {
	var foodConsumptions []*model.FoodConsumption
	err := r.db.NewSelect().Model(&foodConsumptions).Scan(r.ctx)
	return foodConsumptions, err
}

func (r *FoodConsumptionRepository) FindAllFoodConsumptionForMeal(mealId uuid.UUID) ([]*model.FoodConsumption, error) {
	var foodConsumptions []*model.FoodConsumption
	err := r.db.NewSelect().Model(&foodConsumptions).Where("meal_id = ?", mealId).Scan(r.ctx)
	return foodConsumptions, err
}

func (r *FoodConsumptionRepository) FindById(id uuid.UUID) (*model.FoodConsumption, error) {
	var foodConsumption model.FoodConsumption
	err := r.db.NewSelect().Model(&foodConsumption).Where("id = ?", id).Scan(r.ctx)
	return &foodConsumption, err
}

func (r *FoodConsumptionRepository) Create(foodConsumption *model.FoodConsumption) (sql.Result, error) {
	return r.db.NewInsert().Model(foodConsumption).Exec(r.ctx)
}

func (r *FoodConsumptionRepository) Update(foodConsumption *model.FoodConsumption) (sql.Result, error) {
	return r.db.NewUpdate().Model(foodConsumption).Where("id = ?", foodConsumption.ID).Exec(r.ctx)
}

func (r *FoodConsumptionRepository) Delete(foodConsumption *model.FoodConsumption) (sql.Result, error) {
	return r.db.NewDelete().Model(foodConsumption).Exec(r.ctx)
}

func (r *FoodConsumptionRepository) DeleteAllFoodConsumptionForMeal(mealId uuid.UUID) (sql.Result, error) {
	return r.db.NewDelete().Model(&model.FoodConsumption{}).Where("meal_id = ?", mealId).Exec(r.ctx)
}

func (r *FoodConsumptionRepository) DeleteFoodConsumptionForMeal(mealId uuid.UUID, foodConsumptionId uuid.UUID) (sql.Result, error) {
	return r.db.NewDelete().Model(&model.FoodConsumption{}).Where("meal_id = ?", mealId).Where("id = ?", foodConsumptionId).Exec(r.ctx)
}

func (r *FoodConsumptionRepository) GetKcalSumForMeal(mealId uuid.UUID) (float32, error) {
	var sum float32
	err := r.db.NewSelect().ColumnExpr("SUM(kcal)").Table("food_consumption").Where("meal_id = ?", mealId).Scan(r.ctx, &sum)
	return sum, err
}

func (r *FoodConsumptionRepository) GetCostSumForMeal(mealId uuid.UUID) (float32, error) {
	var sum float32
	err := r.db.NewSelect().ColumnExpr("SUM(cost)").Table("food_consumption").Where("meal_id = ?", mealId).Scan(r.ctx, &sum)
	return sum, err
}

func (r *FoodConsumptionRepository) GetMostConsumedFoodInDateRange(startRange time.Time, endRange time.Time, userId string) (*dto.MostConsumedFoodDto, error) {
	var mostConsumedFoodDto dto.MostConsumedFoodDto
	endRange = setEndOfTheDay(endRange)

	query := "SELECT food_id as foodId, food_name AS foodName, SUM(quantity_used_std) AS quantityUsedStd, SUM(quantity_used) AS quantityUsed, unit  FROM food_consumption where meal_id in (select id from meal where user_id = ? and date >= ? and date <= ?) group by food_id, food_name, unit order by quantityUsedStd desc limit 1"
	queryResult, err := r.db.Query(query, userId, startRange, endRange)
	if err != nil {
		return nil, err
	}
	queryResult.Next()
	err = queryResult.Scan(&mostConsumedFoodDto.FoodId, &mostConsumedFoodDto.FoodName, &mostConsumedFoodDto.QuantityUsedStd, &mostConsumedFoodDto.QuantityUsed, &mostConsumedFoodDto.Unit)
	queryResult.Close()
	if err != nil {
		return nil, err
	}
	return &mostConsumedFoodDto, nil
}
