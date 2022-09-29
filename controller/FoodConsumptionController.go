package controller

import (
	"food-track-be/model/dto"
	"food-track-be/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FoodConsumptionController struct {
	foodConsumptionService *service.FoodConsumptionService
}

func NewFoodConsumptionController(foodConsumptionService *service.FoodConsumptionService) *FoodConsumptionController {
	return &FoodConsumptionController{foodConsumptionService: foodConsumptionService}
}

func (s *FoodConsumptionController) FindAllConsumptionForMeal(c *gin.Context) {
	mealId, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	foodConsumptionDtos, err := s.foodConsumptionService.FindAllFoodConsumptionForMeal(mealId)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	c.JSON(200, dto.BaseResponse[[]*dto.FoodConsumptionDto]{
		Body: foodConsumptionDtos,
	})
}

func (s *FoodConsumptionController) AddFoodConsumption(c *gin.Context) {
	mealId, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	var foodConsumptionDto dto.FoodConsumptionDto
	err = c.BindJSON(&foodConsumptionDto)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	foodConsumptionDto, err = s.foodConsumptionService.CreateFoodConsumptionForMeal(mealId, foodConsumptionDto)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	c.JSON(200, dto.BaseResponse[dto.FoodConsumptionDto]{
		Body: foodConsumptionDto,
	})
}

func (s *FoodConsumptionController) UpdateFoodConsumption(c *gin.Context) {
	mealId, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	var foodConsumptionDto dto.FoodConsumptionDto
	c.BindJSON(&foodConsumptionDto)
	foodConsumptionDto, err = s.foodConsumptionService.UpdateFoodConsumptionForMeal(mealId, foodConsumptionDto)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	c.JSON(200, dto.BaseResponse[dto.FoodConsumptionDto]{
		Body: foodConsumptionDto,
	})
}

func (s *FoodConsumptionController) DeleteFoodConsumption(c *gin.Context) {
	mealId, err := uuid.Parse(c.Param("mealId"))
	foodConsumptionId, err := uuid.Parse(c.Param("foodConsumptionId"))
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	err = s.foodConsumptionService.DeleteFoodConsumptionForMeal(mealId, foodConsumptionId)
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	c.JSON(200, dto.BaseResponse[any]{
		Body: err != nil,
	})
}
