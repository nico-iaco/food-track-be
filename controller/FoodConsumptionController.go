package controller

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"food-track-be/model/dto"
	"food-track-be/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

type FoodConsumptionController struct {
	foodConsumptionService *service.FoodConsumptionService
	firebaseApp            *firebase.App
}

func NewFoodConsumptionController(foodConsumptionService *service.FoodConsumptionService, fa *firebase.App) *FoodConsumptionController {
	return &FoodConsumptionController{foodConsumptionService: foodConsumptionService, firebaseApp: fa}
}

func (s *FoodConsumptionController) FindAllConsumptionForMeal(c *gin.Context) {
	_, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	mealId, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	foodConsumptionDtos, err := s.foodConsumptionService.FindAllFoodConsumptionForMeal(mealId)
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	c.JSON(200, dto.BaseResponse[[]*dto.FoodConsumptionDto]{
		Body: foodConsumptionDtos,
	})
}

func (s *FoodConsumptionController) AddFoodConsumption(c *gin.Context) {
	mealId, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	userId, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	var foodConsumptionDto dto.FoodConsumptionDto
	err = c.BindJSON(&foodConsumptionDto)
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	foodConsumptionDto, err = s.foodConsumptionService.CreateFoodConsumptionForMeal(mealId, foodConsumptionDto, userId)
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	c.JSON(200, dto.BaseResponse[dto.FoodConsumptionDto]{
		Body: foodConsumptionDto,
	})
}

func (s *FoodConsumptionController) UpdateFoodConsumption(c *gin.Context) {
	mealId, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	userId, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	var foodConsumptionDto dto.FoodConsumptionDto
	c.BindJSON(&foodConsumptionDto)
	foodConsumptionDto, err = s.foodConsumptionService.UpdateFoodConsumptionForMeal(mealId, foodConsumptionDto, userId)
	if err != nil {
		s.abortWithMessage(c, err.Error())
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
		s.abortWithMessage(c, err.Error())
		return
	}
	userId, err := s.validateTokenAndGetUserId(c.GetHeader("Authorization"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	err = s.foodConsumptionService.DeleteFoodConsumptionForMeal(mealId, foodConsumptionId, userId)
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	c.JSON(200, dto.BaseResponse[any]{
		Body: err != nil,
	})
}

func (s *FoodConsumptionController) abortWithMessage(c *gin.Context, message string) {
	c.AbortWithStatusJSON(200, dto.BaseResponse[any]{
		ErrorMessage: message,
	})
}

func (s *FoodConsumptionController) validateTokenAndGetUserId(authHeader string) (string, error) {
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
