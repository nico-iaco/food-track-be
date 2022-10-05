package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type FoodDetailService struct {
	baseUrl string
}

func NewFoodDetailService() *FoodDetailService {
	return &FoodDetailService{baseUrl: os.Getenv("FOOD_DETAIL_BASE_URL")}
}

func (s FoodDetailService) getCall(url string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

func (s FoodDetailService) GetKcalsForFoodConsumed(barcode string, quantity float64) (float64, error) {
	var response float64
	responseData, err := s.getCall(s.baseUrl + "/food/" + barcode + "/kcal?quantity=" + fmt.Sprintf("%f", quantity) + "&unit=g")
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return 0, err
	}
	return response, nil
}
