package controller

import (
	"food-track-be/model/dto"
	"food-track-be/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MealController struct {
	mealService *service.MealService
}

func NewMealController(mealService *service.MealService) *MealController {
	return &MealController{mealService: mealService}
}

func (s *MealController) FindAllMeals(c *gin.Context) {
	mealDtos, err := s.mealService.FindAll()
	response := dto.BaseResponse[[]*dto.MealDto]{
		Body:         mealDtos,
		ErrorMessage: err.Error(),
	}
	c.JSON(200, response)
}

func (s *MealController) FindMealById(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	mealDto, err := s.mealService.FindById(id)
	response := dto.BaseResponse[dto.MealDto]{
		Body:         mealDto,
		ErrorMessage: err.Error(),
	}
	c.JSON(200, response)
}

func (s *MealController) CreateMeal(c *gin.Context) {
	var mealDto dto.MealDto
	c.BindJSON(&mealDto)
	mealDto, err := s.mealService.Create(mealDto)
	response := dto.BaseResponse[dto.MealDto]{
		Body:         mealDto,
		ErrorMessage: err.Error(),
	}
	c.JSON(200, response)
}

func (s *MealController) UpdateMeal(c *gin.Context) {
	var mealDto dto.MealDto
	c.BindJSON(&mealDto)
	mealDto, err := s.mealService.Update(mealDto)
	response := dto.BaseResponse[dto.MealDto]{
		Body:         mealDto,
		ErrorMessage: err.Error(),
	}
	c.JSON(200, response)
}

func (s *MealController) DeleteMeal(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	err := s.mealService.Delete(id)
	response := dto.BaseResponse[bool]{
		Body:         err == nil,
		ErrorMessage: err.Error(),
	}
	c.JSON(200, response)
}
