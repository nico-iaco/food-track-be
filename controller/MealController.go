package controller

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"food-track-be/model/dto"
	"food-track-be/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
	"time"
)

type MealController struct {
	mealService *service.MealService
	firebaseApp *firebase.App
}

func NewMealController(mealService *service.MealService, fa *firebase.App) *MealController {
	return &MealController{mealService: mealService, firebaseApp: fa}
}

func (s *MealController) FindAllMeals(c *gin.Context) {
	var mealDtos = make([]dto.MealDto, 0)
	startRangeParam := c.Query("startRange")
	endRangeParam := c.Query("endRange")
	userId, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	if startRangeParam != "" && endRangeParam != "" {
		startRange, err := time.Parse("02-01-2006", startRangeParam)
		if err != nil {
			s.abortWithMessage(c, err.Error())
			return
		}
		endRange, err := time.Parse("02-01-2006", endRangeParam)
		if err != nil {
			endRange = startRange
		}
		mealDtos, err = s.mealService.FindAllInDateRange(startRange, endRange, userId)
		if err != nil {
			s.abortWithMessage(c, err.Error())
			return
		}
	} else {
		var err error
		mealDtos, err = s.mealService.FindAll(userId)
		if err != nil {
			s.abortWithMessage(c, err.Error())
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
	userId, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	mealDto, err := s.mealService.FindById(id, userId)
	if err != nil {
		s.abortWithMessage(c, err.Error())
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
		s.abortWithMessage(c, err.Error())
		return
	}
	userId, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	mealDto.UserId = userId
	mealDto, err = s.mealService.Create(mealDto)
	if err != nil {
		s.abortWithMessage(c, err.Error())
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
	userId, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	mealDto, err = s.mealService.Update(mealDto, userId)
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	response := dto.BaseResponse[dto.MealDto]{
		Body: mealDto,
	}
	c.JSON(200, response)
}

func (s *MealController) DeleteMeal(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("mealId"))
	userId, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	err = s.mealService.Delete(id, userId)
	if err != nil {
		s.abortWithMessage(c, err.Error())
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
	userId, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	if startRangeParam != "" && endRangeParam != "" {
		startRange, err := time.Parse("02-01-2006", startRangeParam)
		if err != nil {
			s.abortWithMessage(c, err.Error())
			return
		}
		endRange, err := time.Parse("02-01-2006", endRangeParam)
		if err != nil {
			endRange = startRange
		}
		mealStatisticsDto, err = s.mealService.GetMealsStatistics(startRange, endRange, userId)
		if err != nil {
			s.abortWithMessage(c, err.Error())
			return
		}
	} else {
		startRange := time.Now().AddDate(0, 0, -7)
		endRange := time.Now()
		var err error
		mealStatisticsDto, err = s.mealService.GetMealsStatistics(startRange, endRange, userId)
		if err != nil {
			s.abortWithMessage(c, err.Error())
			return
		}
	}

	response := dto.BaseResponse[dto.MealStatisticsDto]{
		Body: mealStatisticsDto,
	}
	c.JSON(200, response)
}

func (s *MealController) abortWithMessage(c *gin.Context, message string) {
	c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
		ErrorMessage: message,
	})
}

func (s *MealController) validateTokenAndGetUserId(authHeader string) (string, error) {
	auth, err := s.firebaseApp.Auth(context.Background())
	filteredToken := strings.Replace(authHeader, "Bearer ", "", 1)
	if err != nil {
		return "", err
	}
	token, err := auth.VerifyIDToken(context.Background(), filteredToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}
