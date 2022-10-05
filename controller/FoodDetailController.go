package controller

import (
	"food-track-be/model/dto"
	"food-track-be/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FoodDetailController struct {
	foodDetailService *service.FoodDetailService
}

func NewFoodDetailController(f *service.FoodDetailService) *FoodDetailController {
	return &FoodDetailController{foodDetailService: f}
}

func (s *FoodDetailController) GetFoodKcals(c *gin.Context) {
	var baseResponse dto.BaseResponse[float64]
	barcode := c.Param("barcode")
	quantity := c.Query("quantity")
	parsedQuantity, err := strconv.ParseFloat(quantity, 64)
	kcals, err := s.foodDetailService.GetKcalsForFoodConsumed(barcode, parsedQuantity)
	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	baseResponse.Body = kcals
	c.JSON(200, baseResponse)
}
