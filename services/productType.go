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

func (s *ProductTypeService) Update(id, name string) (models.ProductType, error) {
	var productType models.ProductType
	if err := database.DB.Where("id = ?", id).First(&productType).Error; err != nil {
		return productType, err
	}
	productType.Name = name
	err := database.DB.Save(&productType).Error
	return productType, err
}

func (s *ProductTypeService) Delete(id string) error {
	return database.DB.Where("id = ?", id).Delete(&models.ProductType{}).Error
}