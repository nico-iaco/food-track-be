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

// FindAll retrieves all food consumption records from the database.
func (r *FoodConsumptionRepository) FindAll() ([]*model.FoodConsumption, error) {
	// Initialize a slice to hold the retrieved food consumption records.
	var foodConsumptions []*model.FoodConsumption

	// Execute a SELECT statement to retrieve all food consumption records from the database.
	// The result will be stored in the foodConsumptions slice.
	err := r.db.NewSelect().Model(&foodConsumptions).Scan(r.ctx)

	// Return the retrieved food consumption records and any errors.
	return foodConsumptions, err
}

// FindAllFoodConsumptionForMeal retrieves all food consumption records for a specific meal from the database.
func (r *FoodConsumptionRepository) FindAllFoodConsumptionForMeal(mealId uuid.UUID) ([]*model.FoodConsumption, error) {
	// Initialize a slice to hold the retrieved food consumption records.
	var foodConsumptions []*model.FoodConsumption

	// Execute a SELECT statement to retrieve all food consumption records for the specified meal from the database.
	// The result will be stored in the foodConsumptions slice.
	err := r.db.NewSelect().Model(&foodConsumptions).Where("meal_id = ?", mealId).Scan(r.ctx)

	// Return the slice and any error that may have occurred.
	return foodConsumptions, err
}

// FindById retrieves a single food consumption record from the database based on its ID.
func (r *FoodConsumptionRepository) FindById(id uuid.UUID) (*model.FoodConsumption, error) {
	// Initialize a foodConsumption struct to hold the retrieved food consumption record.
	var foodConsumption model.FoodConsumption

	// Execute a SELECT statement to retrieve the food consumption record with the specified ID from the database.
	// The result will be stored in the foodConsumption struct.
	err := r.db.NewSelect().Model(&foodConsumption).Where("id = ?", id).Scan(r.ctx)

	// Return a pointer to the foodConsumption struct and any error that may have occurred.
	return &foodConsumption, err
}

// Create inserts a new food consumption record into the database.
func (r *FoodConsumptionRepository) Create(foodConsumption *model.FoodConsumption) (sql.Result, error) {
	// Execute an INSERT statement to insert the foodConsumption struct as a new row in the database.
	// The result will be stored in a sql.Result value.
	return r.db.NewInsert().Model(foodConsumption).Exec(r.ctx)
}

func (r *FoodConsumptionRepository) Update(foodConsumption *model.FoodConsumption) (sql.Result, error) {
	// Execute an UPDATE statement to update the food consumption record with the specified ID in the database.
	// The result will be stored in a sql.Result value.
	return r.db.NewUpdate().Model(foodConsumption).Where("id = ?", foodConsumption.ID).Exec(r.ctx)
}

// Delete deletes an existing food consumption record from the database.
func (r *FoodConsumptionRepository) Delete(foodConsumption *model.FoodConsumption) (sql.Result, error) {
	// Execute a DELETE statement to delete the food consumption record with the specified ID from the database.
	// The result will be stored in a sql.Result value.
	return r.db.NewDelete().Model(foodConsumption).Exec(r.ctx)
}

// DeleteAllFoodConsumptionForMeal deletes all food consumption records for a particular meal from the database.
func (r *FoodConsumptionRepository) DeleteAllFoodConsumptionForMeal(mealId uuid.UUID) (sql.Result, error) {
	// Execute a DELETE statement to delete all food consumption records with the specified meal ID from the database.
	// The result will be stored in a sql.Result value.
	return r.db.NewDelete().Model(&model.FoodConsumption{}).Where("meal_id = ?", mealId).Exec(r.ctx)
}

// DeleteFoodConsumptionForMeal deletes a specific food consumption record for a particular meal from the database.
func (r *FoodConsumptionRepository) DeleteFoodConsumptionForMeal(mealId uuid.UUID, foodConsumptionId uuid.UUID) (sql.Result, error) {
	// Execute a DELETE statement to delete the food consumption record with the specified IDs from the database.
	// The result will be stored in a sql.Result value.
	return r.db.NewDelete().Model(&model.FoodConsumption{}).Where("meal_id = ?", mealId).Where("id = ?", foodConsumptionId).Exec(r.ctx)
}

// GetKcalSumForMeal retrieves the sum of the "kcal" column for all food consumption records belonging to a particular meal from the database.
func (r *FoodConsumptionRepository) GetKcalSumForMeal(mealId uuid.UUID) (float32, error) {
	// Declare a variable to store the sum of the "kcal" column.
	var sum float32
	// Execute a SELECT statement to retrieve the sum of the "kcal" column for all food consumption records with the specified meal ID.
	// The sum will be stored in the "sum" variable.
	err := r.db.NewSelect().ColumnExpr("SUM(kcal)").Table("food_consumption").Where("meal_id = ?", mealId).Scan(r.ctx, &sum)
	// Return the sum and any error that occurred.
	return sum, err
}

// GetCostSumForMeal retrieves the sum of the "cost" column for all food consumption records belonging to a particular meal from the database.
func (r *FoodConsumptionRepository) GetCostSumForMeal(mealId uuid.UUID) (float32, error) {
	// Declare a variable to store the sum of the "cost" column.
	var sum float32
	// Execute a SELECT statement to retrieve the sum of the "cost" column for all food consumption records with the specified meal ID.
	// The sum will be stored in the "sum" variable.
	err := r.db.NewSelect().ColumnExpr("SUM(cost)").Table("food_consumption").Where("meal_id = ?", mealId).Scan(r.ctx, &sum)
	// Return the sum and any error that occurred.
	return sum, err
}

// GetMostConsumedFoodInDateRange retrieves the food that was consumed the most (by standard quantity used) in a given date range for a particular user from the database.
func (r *FoodConsumptionRepository) GetMostConsumedFoodInDateRange(startRange time.Time, endRange time.Time, userId string) (*dto.MostConsumedFoodDto, error) {
	// Declare a variable to store the most consumed food.
	var mostConsumedFoodDto dto.MostConsumedFoodDto

	// Set the end of the day for the end range.
	endRange = setEndOfTheDay(endRange)

	// Define the SELECT statement to retrieve the most consumed food.
	query := "SELECT food_id as foodId, food_name AS foodName, SUM(quantity_used_std) AS quantityUsedStd, SUM(quantity_used) AS quantityUsed, unit  FROM food_consumption where meal_id in (select id from meal where user_id = ? and date >= ? and date <= ?) group by food_id, food_name, unit order by quantityUsedStd desc limit 1"
	// Execute the SELECT statement and scan the result into the "mostConsumedFoodDto" variable.
	err := r.db.QueryRow(query, userId, startRange, endRange).Scan(&mostConsumedFoodDto.FoodId, &mostConsumedFoodDto.FoodName, &mostConsumedFoodDto.QuantityUsedStd, &mostConsumedFoodDto.QuantityUsed, &mostConsumedFoodDto.Unit)
	// Return the most consumed food or any error that occurred.
	if err != nil {
		return nil, err
	}
	return &mostConsumedFoodDto, nil
}
