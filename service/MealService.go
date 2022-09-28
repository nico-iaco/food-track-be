package service

import (
	"food-track-be/model"
	"food-track-be/model/dto"
	"food-track-be/repository"
	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type MealService struct {
	repository *repository.MealRepository
}

func NewMealService(repository *repository.MealRepository) *MealService {
	return &MealService{repository: repository}
}

func (s *MealService) FindAll() ([]*dto.MealDto, error) {
	var mealsDto []*dto.MealDto
	meals, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}
	mappedField := smapping.MapFields(meals)
	err = smapping.FillStruct(&mealsDto, mappedField)
	if err != nil {
		return []*dto.MealDto{}, err
	}
	return mealsDto, nil
}

func (s *MealService) FindById(id uuid.UUID) (dto.MealDto, error) {
	mealDto := dto.MealDto{}
	meal, err := s.repository.FindById(id)
	if err != nil {
		return mealDto, err
	}
	mappedField := smapping.MapFields(meal)
	err = smapping.FillStruct(&mealDto, mappedField)
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

func (s *MealService) Update(mealDto dto.MealDto) (dto.MealDto, error) {
	meal := model.Meal{}
	mappedField := smapping.MapFields(&mealDto)
	err := smapping.FillStruct(&meal, mappedField)
	if err != nil {
		return mealDto, err
	}
	_, err = s.repository.Update(&meal)
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

func (s *MealService) Delete(mealId uuid.UUID) error {
	meal, err := s.repository.FindById(mealId)
	if err != nil {
		return err
	}
	_, err = s.repository.Delete(meal)
	if err != nil {
		return err
	}
	return nil
}
