package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Meal struct {
	bun.BaseModel    `bun:"table:meal,alias:m"`
	ID               uuid.UUID          `bun:"type:uuid,nullzero,pk"`
	UserId           string             `bun:"type:varchar(255),notnull"`
	Name             string             `bun:"type:varchar(255),notnull"`
	Description      string             `bun:"type:varchar(255),nullzero"`
	MealType         MealType           `bun:"type:varchar(30),notnull"`
	Date             time.Time          `bun:"type:timestamp,notnull"`
	FoodConsumptions []*FoodConsumption `bun:"rel:has-many,join:id=meal_id"`
	//FoodTypes        []FoodType         `bun:"type:varchar(255)[]"`
}

type MealType string

const (
	Breakfast MealType = "breakfast"
	Lunch     MealType = "lunch"
	Dinner    MealType = "dinner"
	Others    MealType = "others"
)

type FoodType string

const (
	Sweet    FoodType = "sweet"
	Sour     FoodType = "sour"
	Salty    FoodType = "salty"
	Bitter   FoodType = "bitter"
	Beverage FoodType = "beverage"
)

/*
DDL for table meal
create table meal (
id uuid primary key,
user_id varchar(255) not null,
name varchar(255) not null,
description varchar(255),
meal_type varchar(255) not null,
date date not null
);
*/
