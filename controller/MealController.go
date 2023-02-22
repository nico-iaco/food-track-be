package controller

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"food-track-be/model/dto"
	"food-track-be/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
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

// FindAllMeals godoc
//	@Summary		Get all meals
//	@Description	get all the meals which satisfies the query parameters (startRange, endRange) or all meals if no query parameters are provided
//	@Tags			meal
//	@Produce		json
//	@Param			startRange	query		string	false	"Start date of the range"
//	@Param			endRange	query		string	false	"End date of the range"
//	@Success		200			{object}	dto.BaseResponse[[]dto.MealDto]
//	@Router			/ [get]
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

// FindMealById godoc
//	@Summary		Get meal
//	@Description	get the meal with the provided id
//	@Tags			meal
//	@Produce		json
//	@Param			mealId	path		string	true	"Meal ID"
//	@Success		200		{object}	dto.BaseResponse[dto.MealDto]
//	@Router			/{mealId}/ [get]
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

// CreateMeal godoc
//	@Summary		Create meal
//	@Description	create a new meal
//	@Tags			meal
//	@Accept			json
//	@Produce		json
//	@Param			mealDto	body		dto.MealDto	true	"Meal to create"
//	@Success		200		{object}	dto.BaseResponse[dto.MealDto]
//	@Router			/ [post]
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

// UpdateMeal godoc
//	@Summary		Update meal
//	@Description	update the meal with the provided id
//	@Tags			meal
//	@Accept			json
//	@Produce		json
//	@Param			mealId	path		string		true	"Meal ID"
//	@Param			mealDto	body		dto.MealDto	true	"Meal to create"
//	@Success		200		{object}	dto.BaseResponse[dto.MealDto]
//	@Router			/{mealId}/ [patch]
func (s *MealController) UpdateMeal(c *gin.Context) {
	var mealDto dto.MealDto
	id, _ := uuid.Parse(c.Param("mealId"))
	err := c.BindJSON(&mealDto)
	userId, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	mealDto.ID = id
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

// DeleteMeal godoc
//	@Summary		Delete meal
//	@Description	delete the meal with the provided id
//	@Tags			meal
//	@Accept			json
//	@Produce		json
//	@Param			mealId	path		string	true	"Meal ID"
//	@Success		200		{object}	dto.BaseResponse[bool]
//	@Router			/{mealId}/ [delete]
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

// GetMealStatistics godoc
//	@Summary		Get meal statistics
//	@Description	get the meal statistics for the provided date range (default is the past week)
//	@Tags			meal
//	@Produce		json
//	@Param			startRange	query		string	false	"Start date of the range"
//	@Param			endRange	query		string	false	"End date of the range"
//	@Success		200			{object}	dto.BaseResponse[dto.MealStatisticsDto]
//	@Router			/statistics/ [get]
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
	log.Println(message)
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
