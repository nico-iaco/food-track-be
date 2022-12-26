package service

import (
	"food-track-be/model"
	"food-track-be/model/dto"
	"food-track-be/repository"
	"github.com/google/uuid"
	"github.com/mashingan/smapping"
	"time"
)

type MealService struct {
	repository             *repository.MealRepository
	foodConsumptionService *FoodConsumptionService
}

func NewMealService(repository *repository.MealRepository, service *FoodConsumptionService) *MealService {
	return &MealService{repository: repository, foodConsumptionService: service}
}

func (s *MealService) FindAll(userId string) ([]dto.MealDto, error) {
	var mealsDto []dto.MealDto
	meals, err := s.repository.FindAll(userId)
	if err != nil {
		return nil, err
	}
	for _, meal := range meals {
		mealDto := dto.MealDto{}
		mappedField := smapping.MapFields(&meal)
		err = smapping.FillStruct(&mealDto, mappedField)
		if err != nil {
			return nil, err
		}
		mealDto.Kcal, err = s.foodConsumptionService.GetKcalSumForMeal(meal.ID)
		if err != nil {
			return nil, err
		}
		mealDto.Cost, err = s.foodConsumptionService.GetCostSumForMeal(meal.ID)
		if err != nil {
			return nil, err
		}
		mealsDto = append(mealsDto, mealDto)
	}
	if err != nil {
		return []dto.MealDto{}, err
	}
	return mealsDto, nil
}

func (s *MealService) FindAllInDateRange(startRange time.Time, endRange time.Time, userId string) ([]dto.MealDto, error) {
	var mealsDto []dto.MealDto
	meals, err := s.repository.GetMealInDateRange(startRange, endRange, userId)
	if err != nil {
		return nil, err
	}
	for _, meal := range meals {
		mealDto := dto.MealDto{}
		mappedField := smapping.MapFields(&meal)
		err = smapping.FillStruct(&mealDto, mappedField)
		if err != nil {
			return nil, err
		}
		mealDto.Kcal, err = s.foodConsumptionService.GetKcalSumForMeal(meal.ID)
		if err != nil {
			return nil, err
		}
		mealDto.Cost, err = s.foodConsumptionService.GetCostSumForMeal(meal.ID)
		if err != nil {
			return nil, err
		}
		mealsDto = append(mealsDto, mealDto)
	}
	if err != nil {
		return []dto.MealDto{}, err
	}
	return mealsDto, nil
}

func (s *MealService) FindById(id uuid.UUID, userId string) (dto.MealDto, error) {
	mealDto := dto.MealDto{}
	meal, err := s.repository.FindByIdAndUserId(id, userId)
	if err != nil {
		return mealDto, err
	}
	mappedField := smapping.MapFields(meal)
	err = smapping.FillStruct(&mealDto, mappedField)
	if err != nil {
		return mealDto, err
	}
	mealDto.Kcal, err = s.foodConsumptionService.GetKcalSumForMeal(meal.ID)
	if err != nil {
		return mealDto, err
	}
	mealDto.Cost, err = s.foodConsumptionService.GetCostSumForMeal(meal.ID)
	if err != nil {
		return mealDto, err
	}
	return mealDto, nil
}

func (s *MealService) Create(mealDto dto.MealDto) (dto.MealDto, error) {
	meal := model.Meal{}
	mappedField := smapping.MapFields(&mealDto)
	err := smapping.FillStruct(&meal, mappedField)
	if err != nil {
		return mealDto, err
	}
	meal.ID = uuid.New()
	_, err = s.repository.Create(&meal)
	if err != nil {
		return dto.MealDto{}, err
	}
	mappedField = smapping.MapFields(&meal)
	err = smapping.FillStruct(&mealDto, mappedField)
	if err != nil {
		return mealDto, err
	}
	return mealDto, nil
}

func (s *MealService) Update(mealDto dto.MealDto, userId string) (dto.MealDto, error) {
	meal := model.Meal{}
	mappedField := smapping.MapFields(&mealDto)
	err := smapping.FillStruct(&meal, mappedField)
	if err != nil {
		return mealDto, err
	}
	_, err = s.repository.Update(&meal, userId)
	if err != nil {
		return dto.MealDto{}, err
	}
	mappedField = smapping.MapFields(&meal)
	err = smapping.FillStruct(&mealDto, mappedField)
	if err != nil {
		return mealDto, err
	}
	mealDto.Kcal, err = s.foodConsumptionService.GetKcalSumForMeal(meal.ID)
	if err != nil {
		return mealDto, err
	}
	mealDto.Cost, err = s.foodConsumptionService.GetCostSumForMeal(meal.ID)
	if err != nil {
		return mealDto, err
	}
	return mealDto, nil
}

func (s *MealService) Delete(mealId uuid.UUID, userId string) error {
	meal, err := s.repository.FindByIdAndUserId(mealId, userId)
	if err != nil {
		return err
	}
	_, err = s.repository.Delete(meal, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *MealService) GetMealsStatistics(startRange time.Time, endRange time.Time, userId string) (dto.MealStatisticsDto, error) {
	var mealStatisticsDto dto.MealStatisticsDto
	avgKcal, err := s.repository.GetAverageKcalEatenInDateRange(startRange, endRange, userId)
	if err != nil {
		return dto.MealStatisticsDto{}, err
	}
	mealStatisticsDto.AverageWeekCalories = avgKcal

	avgKcalPerMealType, err := s.repository.GetAverageKcalEatenInDateRangePerMealType(startRange, endRange, userId)
	if err != nil {
		return dto.MealStatisticsDto{}, err
	}
	mealStatisticsDto.AverageWeekCaloriesPerMealType = avgKcalPerMealType

	avgCost, err := s.repository.GetAverageFoodCostInDateRange(startRange, endRange, userId)
	if err != nil {
		return dto.MealStatisticsDto{}, err
	}
	mealStatisticsDto.AverageWeekFoodCost = avgCost

	sumFoodCost, err := s.repository.GetSumFoodCostInDateRange(startRange, endRange, userId)
	if err != nil {
		return dto.MealStatisticsDto{}, err
	}
	mealStatisticsDto.SumWeekFoodCost = sumFoodCost

	mostConsumedFood, err := s.foodConsumptionService.GetMostConsumedFoodInDateRange(startRange, endRange, userId)
	if err != nil {
		return dto.MealStatisticsDto{}, err
	}
	mealStatisticsDto.MostConsumedFood = *mostConsumedFood

	return mealStatisticsDto, nil
}
