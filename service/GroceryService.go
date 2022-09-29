package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"food-track-be/model/dto"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
)

type GroceryService struct {
	baseUrl string
}

func NewGroceryService() *GroceryService {
	return &GroceryService{baseUrl: os.Getenv("GROCERY_BASE_URL")}
}

func getCall(url string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

func patchCall(url string, body any) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	request, err := http.NewRequest(http.MethodPatch, url, &buf)
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

func (s *GroceryService) GetAllAvailableFood() ([]*dto.FoodAvailableDto, error) {
	var response dto.BaseResponse[[]*dto.FoodAvailableDto]
	responseData, err := getCall(s.baseUrl + "/item/")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (s *GroceryService) GetAvailableTransactionForFood(foodId uuid.UUID) ([]*dto.FoodTransactionDto, error) {
	var response dto.BaseResponse[[]*dto.FoodTransactionDto]
	responseData, err := getCall(s.baseUrl + "/item/" + foodId.String() + "/transaction")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (s *GroceryService) GetTransactionDetail(foodId uuid.UUID, transactionId uuid.UUID) (dto.FoodTransactionDto, error) {
	var response dto.BaseResponse[dto.FoodTransactionDto]
	responseData, err := getCall(s.baseUrl + "/item/" + foodId.String() + "/transaction/" + transactionId.String())
	if err != nil {
		return dto.FoodTransactionDto{}, err
	}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return dto.FoodTransactionDto{}, err
	}
	return response.Body, nil
}

func (s *GroceryService) UpdateFoodTransaction(foodId uuid.UUID, foodTransactionDto dto.FoodTransactionDto) (dto.FoodTransactionDto, error) {
	var response dto.BaseResponse[dto.FoodTransactionDto]
	result, err := patchCall(s.baseUrl+"/item/"+foodId.String()+"/transaction/", foodTransactionDto)
	if err != nil {
		return dto.FoodTransactionDto{}, err
	}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return dto.FoodTransactionDto{}, err
	}
	if response.ErrorMessage != "" {
		return dto.FoodTransactionDto{}, errors.New(response.ErrorMessage)
	}
	return response.Body, nil
}
