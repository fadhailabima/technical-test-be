package services

import (
	"technical-test-backend/database"
	"technical-test-backend/models"

	"github.com/google/uuid"
)

type ProductService struct{}

type CreateProductInput struct {
	Name          string  `json:"name" binding:"required"`
	Stock         int     `json:"stock" binding:"required,min=0"`
	Price         float64 `json:"price" binding:"required,min=1"`
	ProductTypeID string  `json:"product_type_id" binding:"required"`
}

func (s *ProductService) Create(input CreateProductInput) (models.Product, error) {
	typeUUID, _ := uuid.Parse(input.ProductTypeID)
	product := models.Product{
		Name: input.Name, Stock: input.Stock, Price: input.Price, ProductTypeID: typeUUID,
	}
	err := database.DB.Create(&product).Error
	return product, err
}

func (s *ProductService) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := database.DB.Preload("ProductType").Find(&products).Error
	return products, err
}

func (s *ProductService) Delete(id string) error {
	return database.DB.Delete(&models.Product{}, "id = ?", id).Error
}
// Update Product - Admin dapat update produk master
type UpdateProductInput struct {
	Name          *string  `json:"name"`
	Stock         *int     `json:"stock" binding:"omitempty,min=0"`
	Price         *float64 `json:"price" binding:"omitempty,min=1"`
	ProductTypeID *string  `json:"product_type_id"`
}

func (s *ProductService) Update(id string, input UpdateProductInput) (models.Product, error) {
	var product models.Product
	
	// Check if product exists
	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		return product, err
	}

	// Update only provided fields
	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Stock != nil {
		updates["stock"] = *input.Stock
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}
	if input.ProductTypeID != nil {
		typeUUID, err := uuid.Parse(*input.ProductTypeID)
		if err != nil {
			return product, err
		}
		updates["product_type_id"] = typeUUID
	}

	if err := database.DB.Model(&product).Updates(updates).Error; err != nil {
		return product, err
	}

	// Reload with ProductType
	database.DB.Preload("ProductType").First(&product, "id = ?", id)
	return product, nil
}
