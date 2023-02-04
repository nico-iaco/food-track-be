package service

import (
	"food-track-be/model"
	"food-track-be/model/dto"
	"food-track-be/repository"
	"github.com/google/uuid"
	"github.com/mashingan/smapping"
	"log"
	"time"
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
	for _, foodConsumption := range foodConsumptions {
		foodConsumptionDto, err := s.mapMealConsumptionToDto(foodConsumption)
		if err != nil {
			return nil, err
		}
		foodConsumptionsDto = append(foodConsumptionsDto, &foodConsumptionDto)
	}
	return foodConsumptionsDto, nil
}

func (s FoodConsumptionService) CreateFoodConsumptionForMeal(mealId uuid.UUID, foodConsumptionDto dto.FoodConsumptionDto, token string) (dto.FoodConsumptionDto, error) {
	foodConsumption := model.FoodConsumption{}
	mappedField := smapping.MapFields(&foodConsumptionDto)
	err := smapping.FillStruct(&foodConsumption, mappedField)
	foodConsumption.MealID = mealId
	if err != nil {
		log.Println(err)
		return foodConsumptionDto, err
	}
	foodConsumption.ID = uuid.New()
	var transactionDto dto.FoodTransactionDto
	if foodConsumption.FoodId != uuid.Nil && foodConsumption.TransactionId != uuid.Nil {
		transactionDto, err = s.groceryService.GetTransactionDetail(foodConsumptionDto.FoodId, foodConsumptionDto.TransactionId, token)
		if err != nil {
			log.Println(err)
			return dto.FoodConsumptionDto{}, err
		}
		foodConsumption.Cost = (transactionDto.Price / transactionDto.Quantity) * foodConsumptionDto.QuantityUsed
	}

	_, err = s.repository.Create(&foodConsumption)

	if foodConsumption.FoodId != uuid.Nil && foodConsumption.TransactionId != uuid.Nil {
		if err != nil {
			transactionDto.AvailableQuantity -= foodConsumptionDto.QuantityUsed
			_, err = s.groceryService.UpdateFoodTransaction(foodConsumptionDto.FoodId, transactionDto, token)
			println(err)
			return dto.FoodConsumptionDto{}, err
		}
	}

	foodConsumptionDto, err = s.mapMealConsumptionToDto(&foodConsumption)
	if err != nil {
		log.Println(err)
		return foodConsumptionDto, err
	}
	return foodConsumptionDto, nil
}

func (s FoodConsumptionService) UpdateFoodConsumptionForMeal(mealId uuid.UUID, foodConsumptionDto dto.FoodConsumptionDto, token string) (dto.FoodConsumptionDto, error) {
	foodConsumption := model.FoodConsumption{}
	err := smapping.FillStruct(&foodConsumption, smapping.MapFields(&foodConsumptionDto))
	foodConsumption.MealID = mealId
	if err != nil {
		return foodConsumptionDto, err
	}

	var transactionDto dto.FoodTransactionDto
	if foodConsumption.FoodId != uuid.Nil && foodConsumption.TransactionId != uuid.Nil {
		transactionDto, err := s.groceryService.GetTransactionDetail(foodConsumptionDto.FoodId, foodConsumptionDto.TransactionId, token)
		if err != nil {
			return dto.FoodConsumptionDto{}, err
		}

		foodConsumption.Cost = (transactionDto.Price / transactionDto.Quantity) * foodConsumptionDto.QuantityUsed

	}

	_, err = s.repository.Update(&foodConsumption)

	if foodConsumption.FoodId != uuid.Nil && foodConsumption.TransactionId != uuid.Nil {
		prevConsumption, err := s.repository.FindById(foodConsumptionDto.ID)
		if err != nil {
			return dto.FoodConsumptionDto{}, err
		}
		deltaQuantity := foodConsumptionDto.QuantityUsed - prevConsumption.QuantityUsed
		transactionDto.AvailableQuantity += deltaQuantity
		_, err = s.groceryService.UpdateFoodTransaction(foodConsumptionDto.FoodId, transactionDto, token)
	}

	foodConsumptionDto, err = s.mapMealConsumptionToDto(&foodConsumption)
	if err != nil {
		return foodConsumptionDto, err
	}
	return foodConsumptionDto, nil
}

func (s FoodConsumptionService) DeleteFoodConsumptionForMeal(mealId uuid.UUID, foodConsumptionId uuid.UUID, token string) error {
	foodConsumption, err := s.repository.FindById(foodConsumptionId)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s.repository.DeleteFoodConsumptionForMeal(mealId, foodConsumptionId)

	if foodConsumption.FoodId != uuid.Nil && foodConsumption.TransactionId != uuid.Nil {
		if err != nil {
			transactionDto, err := s.groceryService.GetTransactionDetail(foodConsumption.FoodId, foodConsumption.TransactionId, token)
			if err != nil {
				log.Println(err)
				return err
			}
			transactionDto.AvailableQuantity += foodConsumption.QuantityUsed
			_, err = s.groceryService.UpdateFoodTransaction(foodConsumption.FoodId, transactionDto, token)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}

func (s FoodConsumptionService) GetKcalSumForMeal(mealId uuid.UUID) (float32, error) {
	return s.repository.GetKcalSumForMeal(mealId)
}

func (s FoodConsumptionService) GetCostSumForMeal(mealId uuid.UUID) (float32, error) {
	return s.repository.GetCostSumForMeal(mealId)
}

func (s FoodConsumptionService) GetMostConsumedFoodInDateRange(startDate time.Time, endDate time.Time, userId string) (*dto.MostConsumedFoodDto, error) {
	mostConsumedFood, err := s.repository.GetMostConsumedFoodInDateRange(startDate, endDate, userId)
	if err != nil {
		return &dto.MostConsumedFoodDto{}, err
	}

	return mostConsumedFood, nil
}

func (s FoodConsumptionService) mapMealConsumptionToDto(foodConsumption *model.FoodConsumption) (dto.FoodConsumptionDto, error) {
	foodConsumptionDto := dto.FoodConsumptionDto{}
	err := smapping.FillStruct(&foodConsumptionDto, smapping.MapFields(&foodConsumption))
	return foodConsumptionDto, err
}
