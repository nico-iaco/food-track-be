package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"food-track-be/model/dto"
	"github.com/google/uuid"
	"github.com/sony/gobreaker"
	"io"
	"net/http"
	"os"
	"time"
)

type GroceryService struct {
	baseUrl        string
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewGroceryService() *GroceryService {
	return &GroceryService{
		baseUrl: os.Getenv("GROCERY_BASE_URL"),
		circuitBreaker: gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        "GroceryService",
			MaxRequests: 5,
			Interval:    60 * time.Second,
			Timeout:     5 * time.Second,
		}),
	}
}

func (s *GroceryService) getCall(url string, userId string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("iv-user", userId)
	if err != nil {
		return nil, err
	}
	result, err := s.circuitBreaker.Execute(func() (interface{}, error) {
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, err
		}
		return io.ReadAll(response.Body)
	})
	if err != nil {
		return nil, err
	}

	return result.([]byte), nil
}

func (s *GroceryService) patchCall(url string, body any, userId string) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	request, err := http.NewRequest(http.MethodPatch, url, &buf)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("iv-user", userId)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

func (s *GroceryService) GetAllAvailableFood(userId string) ([]*dto.FoodAvailableDto, error) {
	var response dto.BaseResponse[[]*dto.FoodAvailableDto]
	responseData, err := s.getCall(s.baseUrl+"/api/item/", userId)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (s *GroceryService) GetAvailableTransactionForFood(foodId uuid.UUID, userId string) ([]*dto.FoodTransactionDto, error) {
	var response dto.BaseResponse[[]*dto.FoodTransactionDto]
	responseData, err := s.getCall(s.baseUrl+"/api/item/"+foodId.String()+"/transaction", userId)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (s *GroceryService) GetTransactionDetail(foodId uuid.UUID, transactionId uuid.UUID, userId string) (dto.FoodTransactionDto, error) {
	var response dto.BaseResponse[dto.FoodTransactionDto]
	responseData, err := s.getCall(s.baseUrl+"/api/item/"+foodId.String()+"/transaction/"+transactionId.String(), userId)
	if err != nil {
		return dto.FoodTransactionDto{}, err
	}
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return dto.FoodTransactionDto{}, err
	}
	return response.Body, nil
}

func (s *GroceryService) UpdateFoodTransaction(foodId uuid.UUID, foodTransactionDto dto.FoodTransactionDto, userId string) (dto.FoodTransactionDto, error) {
	var response dto.BaseResponse[dto.FoodTransactionDto]
	result, err := s.patchCall(s.baseUrl+"/api/item/"+foodId.String()+"/transaction/", foodTransactionDto, userId)
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
