package services

import (
	"technical-test-backend/database"
	"technical-test-backend/models"
)

type ProductTypeService struct{}

func (s *ProductTypeService) GetAll() ([]models.ProductType, error) {
	var types []models.ProductType
	err := database.DB.Find(&types).Error
	return types, err
}

func (s *ProductTypeService) Create(name string) (models.ProductType, error) {
	newType := models.ProductType{Name: name}
	err := database.DB.Create(&newType).Error
	return newType, err
}