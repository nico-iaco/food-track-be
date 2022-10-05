package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type FoodConsumption struct {
	bun.BaseModel   `bun:"table:food_consumption,alias:fc"`
	ID              uuid.UUID `bun:"type:uuid,notnull,pk,default:uuid_generate_v4()"`
	MealID          uuid.UUID
	FoodId          uuid.UUID
	TransactionId   uuid.UUID
	FoodName        string
	QuantityUsed    float32
	QuantityUsedStd float32
	Unit            string
	Kcal            float32
}

/*
DDL for table food_consumption (First create Meal table)
create table food_consumption (
id uuid primary key,
meal_id uuid not null,
food_id uuid not null,
transaction_id uuid not null,
food_name varchar(255) not null,
quantity_used float not null,
quantity_used_std float not null,
unit varchar(255) not null,
kcal float not null,
foreign key (meal_id) references meal(id)
);
*/
