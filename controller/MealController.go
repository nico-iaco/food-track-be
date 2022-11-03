package controller

import (
	"food-track-be/model/dto"
	"food-track-be/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type MealController struct {
	mealService *service.MealService
}

func NewMealController(mealService *service.MealService) *MealController {
	return &MealController{mealService: mealService}
}

func (s *MealController) FindAllMeals(c *gin.Context) {
	var mealDtos = make([]dto.MealDto, 0)
	startRangeParam := c.Query("startRange")
	endRangeParam := c.Query("endRange")
	userId := c.GetHeader("iv-user")
	if startRangeParam != "" && endRangeParam != "" {
		startRange, err := time.Parse("02-01-2006", startRangeParam)
		if err != nil {
			c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
				ErrorMessage: err.Error(),
			})
			return
		}
		endRange, err := time.Parse("02-01-2006", endRangeParam)
		if err != nil {
			endRange = startRange
		}
		mealDtos, err = s.mealService.FindAllInDateRange(startRange, endRange, userId)
		if err != nil {
			c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
				ErrorMessage: err.Error(),
			})
			return
		}
	} else {
		var err error
		mealDtos, err = s.mealService.FindAll(userId)
		if err != nil {
			c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
				ErrorMessage: err.Error(),
			})
			return
		}
	}

	response := dto.BaseResponse[[]dto.MealDto]{
		Body: mealDtos,
	}
	c.JSON(200, response)
}

func (s *MealController) FindMealById(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("mealId"))
	userId := c.GetHeader("iv-user")
	mealDto, err := s.mealService.FindById(id, userId)
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
	userId := c.GetHeader("iv-user")
	if err != nil {
		c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
			ErrorMessage: err.Error(),
		})
		return
	}
	mealDto, err = s.mealService.Update(mealDto, userId)
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
	userId := c.GetHeader("iv-user")
	err := s.mealService.Delete(id, userId)
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

func (s *MealController) GetMealStatistics(c *gin.Context) {
	var mealStatisticsDto dto.MealStatisticsDto
	startRangeParam := c.Query("startRange")
	endRangeParam := c.Query("endRange")
	userId := c.GetHeader("iv-user")
	if startRangeParam != "" && endRangeParam != "" {
		startRange, err := time.Parse("02-01-2006", startRangeParam)
		if err != nil {
			c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
				ErrorMessage: err.Error(),
			})
			return
		}
		endRange, err := time.Parse("02-01-2006", endRangeParam)
		if err != nil {
			endRange = startRange
		}
		mealStatisticsDto, err = s.mealService.GetMealsStatistics(startRange, endRange, userId)
		if err != nil {
			c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
				ErrorMessage: err.Error(),
			})
			return
		}
	} else {
		startRange := time.Now().AddDate(0, 0, -7)
		endRange := time.Now()
		var err error
		mealStatisticsDto, err = s.mealService.GetMealsStatistics(startRange, endRange, userId)
		if err != nil {
			c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
				ErrorMessage: err.Error(),
			})
			return
		}
	}

	response := dto.BaseResponse[dto.MealStatisticsDto]{
		Body: mealStatisticsDto,
	}
	c.JSON(200, response)
}
