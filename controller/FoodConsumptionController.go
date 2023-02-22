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
)

type FoodConsumptionController struct {
	foodConsumptionService *service.FoodConsumptionService
	firebaseApp            *firebase.App
}

func NewFoodConsumptionController(foodConsumptionService *service.FoodConsumptionService, fa *firebase.App) *FoodConsumptionController {
	return &FoodConsumptionController{foodConsumptionService: foodConsumptionService, firebaseApp: fa}
}

// FindAllConsumptionForMeal godoc
//	@Summary		Get all consumption for the meal
//	@Description	find all the consumption for the meal by mealId
//	@Tags			food-consumption
//	@Produce		json
//	@Param			mealId	path		string	true	"Meal ID"
//	@Success		200		{object}	dto.BaseResponse[[]dto.FoodConsumptionDto]
//	@Router			/{mealId}/consumption/ [get]
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

// AddFoodConsumption godoc
//	@Summary		Add consumption for the meal
//	@Description	add consumption for the meal by mealId
//	@Tags			food-consumption
//	@Accept			json
//	@Produce		json
//	@Param			mealId				path		string					true	"Meal ID"
//	@Param			foodConsumptionDto	body		dto.FoodConsumptionDto	true	"Food Consumption"
//	@Success		200					{object}	dto.BaseResponse[dto.FoodConsumptionDto]
//	@Router			/{mealId}/consumption/ [post]
func (s *FoodConsumptionController) AddFoodConsumption(c *gin.Context) {
	mealId, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	token := c.GetHeader("Authorization")
	_, err = s.validateTokenAndGetUserId(token)
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
	foodConsumptionDto, err = s.foodConsumptionService.CreateFoodConsumptionForMeal(mealId, foodConsumptionDto, token)
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	c.JSON(200, dto.BaseResponse[dto.FoodConsumptionDto]{
		Body: foodConsumptionDto,
	})
}

// UpdateFoodConsumption godoc
//	@Summary		Update consumption for the meal
//	@Description	update consumption for the meal by mealId
//	@Tags			food-consumption
//	@Accept			json
//	@Produce		json
//	@Param			mealId				path		string					true	"Meal ID"
//	@Param			foodConsumptionDto	body		dto.FoodConsumptionDto	true	"Food Consumption"
//	@Success		200					{object}	dto.BaseResponse[dto.FoodConsumptionDto]
//	@Router			/{mealId}/consumption/ [patch]
func (s *FoodConsumptionController) UpdateFoodConsumption(c *gin.Context) {
	mealId, err := uuid.Parse(c.Param("mealId"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	token := c.GetHeader("Authorization")
	_, err = s.validateTokenAndGetUserId(token)
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	var foodConsumptionDto dto.FoodConsumptionDto
	c.BindJSON(&foodConsumptionDto)
	foodConsumptionDto, err = s.foodConsumptionService.UpdateFoodConsumptionForMeal(mealId, foodConsumptionDto, token)
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	c.JSON(200, dto.BaseResponse[dto.FoodConsumptionDto]{
		Body: foodConsumptionDto,
	})
}

// DeleteFoodConsumption godoc
//	@Summary		Delete consumption for the meal
//	@Description	delete consumption for the meal by mealId
//	@Tags			food-consumption
//	@Accept			json
//	@Produce		json
//	@Param			mealId				path		string	true	"Meal ID"
//	@Param			foodConsumptionId	path		string	true	"Food consumption ID"
//	@Success		200					{object}	dto.BaseResponse[bool]
//	@Router			/{mealId}/consumption/ [delete]
func (s *FoodConsumptionController) DeleteFoodConsumption(c *gin.Context) {
	mealId, err := uuid.Parse(c.Param("mealId"))
	foodConsumptionId, err := uuid.Parse(c.Param("foodConsumptionId"))
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	token := c.GetHeader("Authorization")
	_, err = s.validateTokenAndGetUserId(token)
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	err = s.foodConsumptionService.DeleteFoodConsumptionForMeal(mealId, foodConsumptionId, token)
	if err != nil {
		s.abortWithMessage(c, err.Error())
		return
	}
	c.JSON(200, dto.BaseResponse[bool]{
		Body: err != nil,
	})
}

func (s *FoodConsumptionController) abortWithMessage(c *gin.Context, message string) {
	log.Println(message)
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
