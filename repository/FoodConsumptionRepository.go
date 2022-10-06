package repository

import (
	"context"
	"database/sql"
	"food-track-be/model"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
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
