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
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	response := dto.BaseResponse[[]dto.MealDto]{
		Body: mealDtos,
	}
	c.JSON(200, response)
}

func (s *MealController) FindMealById(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("mealId"))
	mealDto, err := s.mealService.FindById(id)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	response := dto.BaseResponse[dto.MealDto]{
		Body: mealDto,
	}
	c.JSON(200, response)
}

func (s *MealController) CreateMeal(c *gin.Context) {
	var mealDto dto.MealDto
	err := c.BindJSON(&mealDto)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	mealDto, err = s.mealService.Create(mealDto)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	response := dto.BaseResponse[dto.MealDto]{
		Body: mealDto,
	}
	c.JSON(200, response)
}

func (s *MealController) UpdateMeal(c *gin.Context) {
	var mealDto dto.MealDto
	err := c.BindJSON(&mealDto)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	mealDto, err = s.mealService.Update(mealDto)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	response := dto.BaseResponse[dto.MealDto]{
		Body: mealDto,
	}
	c.JSON(200, response)
}

func (s *MealController) DeleteMeal(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("mealId"))
	err := s.mealService.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	response := dto.BaseResponse[bool]{
		Body: err == nil,
	}
	c.JSON(200, response)
}
