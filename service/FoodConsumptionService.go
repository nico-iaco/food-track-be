package service

import (
	"food-track-be/model"
	"food-track-be/model/dto"
	"food-track-be/repository"
	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type FoodConsumptionService struct {
	repository     *repository.FoodConsumptionRepository
	groceryService *GroceryService
}

func NewFoodConsumptionService(repository *repository.FoodConsumptionRepository, groceryService *GroceryService) *FoodConsumptionService {
	return &FoodConsumptionService{repository: repository, groceryService: groceryService}
}

func (s FoodConsumptionService) FindAllFoodConsumptionForMeal(mealId uuid.UUID) ([]*dto.FoodConsumptionDto, error) {
	var foodConsumptionsDto []*dto.FoodConsumptionDto
	foodConsumptions, err := s.repository.FindAllFoodConsumptionForMeal(mealId)
	if err != nil {
		return nil, err
	}
	mappedField := smapping.MapFields(foodConsumptions)
	err = smapping.FillStruct(&foodConsumptionsDto, mappedField)
	if err != nil {
		return []*dto.FoodConsumptionDto{}, err
	}
	return foodConsumptionsDto, nil
}

func (s FoodConsumptionService) CreateFoodConsumptionForMeal(mealId uuid.UUID, foodConsumptionDto dto.FoodConsumptionDto) (dto.FoodConsumptionDto, error) {
	foodConsumption := model.FoodConsumption{}
	mappedField := smapping.MapFields(&foodConsumptionDto)
	err := smapping.FillStruct(&foodConsumption, mappedField)
	foodConsumption.MealID = mealId
	if err != nil {
		return foodConsumptionDto, err
	}

	transactionDto, err := s.groceryService.GetTransactionDetail(foodConsumptionDto.FoodId, foodConsumptionDto.TransactionId)
	if err != nil {
		return dto.FoodConsumptionDto{}, err
	}
	transactionDto.AvailableQuantity -= foodConsumptionDto.QuantityUsed
	_, err = s.groceryService.UpdateFoodTransaction(foodConsumptionDto.TransactionId, transactionDto)
	if err != nil {
		return dto.FoodConsumptionDto{}, err
	}
	_, err = s.repository.Create(&foodConsumption)
	if err != nil {
		transactionDto.AvailableQuantity += foodConsumptionDto.QuantityUsed
		_, err = s.groceryService.UpdateFoodTransaction(foodConsumptionDto.TransactionId, transactionDto)
		if err != nil {
			return dto.FoodConsumptionDto{}, err
		}
		return dto.FoodConsumptionDto{}, err
	}
	mappedField = smapping.MapFields(&foodConsumption)
	err = smapping.FillStruct(&foodConsumptionDto, mappedField)
	if err != nil {
		return foodConsumptionDto, err
	}
	return foodConsumptionDto, nil
}

func (s FoodConsumptionService) UpdateFoodConsumptionForMeal(mealId uuid.UUID, foodConsumptionDto dto.FoodConsumptionDto) (dto.FoodConsumptionDto, error) {
	foodConsumption := model.FoodConsumption{}
	mappedField := smapping.MapFields(&foodConsumptionDto)
	err := smapping.FillStruct(&foodConsumption, mappedField)
	foodConsumption.MealID = mealId
	if err != nil {
		return foodConsumptionDto, err
	}
	prevConsumption, err := s.repository.FindById(foodConsumptionDto.ID)
	if err != nil {
		return dto.FoodConsumptionDto{}, err
	}
	deltaQuantity := foodConsumptionDto.QuantityUsed - prevConsumption.QuantityUsed
	transactionDto, err := s.groceryService.GetTransactionDetail(foodConsumptionDto.FoodId, foodConsumptionDto.TransactionId)
	if err != nil {
		return dto.FoodConsumptionDto{}, err
	}
	transactionDto.AvailableQuantity += deltaQuantity
	_, err = s.groceryService.UpdateFoodTransaction(foodConsumptionDto.TransactionId, transactionDto)
	if err != nil {
		return dto.FoodConsumptionDto{}, err
	}
	_, err = s.repository.Update(&foodConsumption)
	if err != nil {
		transactionDto.AvailableQuantity -= deltaQuantity
		_, err = s.groceryService.UpdateFoodTransaction(foodConsumptionDto.TransactionId, transactionDto)
		if err != nil {
			return dto.FoodConsumptionDto{}, err
		}
		return dto.FoodConsumptionDto{}, err
	}
	mappedField = smapping.MapFields(&foodConsumption)
	err = smapping.FillStruct(&foodConsumptionDto, mappedField)
	if err != nil {
		return foodConsumptionDto, err
	}
	return foodConsumptionDto, nil
}

func (s FoodConsumptionService) DeleteFoodConsumptionForMeal(mealId uuid.UUID, foodConsumptionId uuid.UUID) error {
	foodConsumption, err := s.repository.FindById(foodConsumptionId)
	if err != nil {
		return err
	}
	transactionDto, err := s.groceryService.GetTransactionDetail(foodConsumption.FoodId, foodConsumption.TransactionId)
	if err != nil {
		return err
	}
	transactionDto.AvailableQuantity += foodConsumption.QuantityUsed
	_, err = s.groceryService.UpdateFoodTransaction(foodConsumption.TransactionId, transactionDto)
	if err != nil {
		return err
	}
	_, err = s.repository.DeleteFoodConsumptionForMeal(mealId, foodConsumptionId)
	if err != nil {
		transactionDto.AvailableQuantity -= foodConsumption.QuantityUsed
		_, err = s.groceryService.UpdateFoodTransaction(foodConsumption.TransactionId, transactionDto)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}