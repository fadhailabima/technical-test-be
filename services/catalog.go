package services

import (
	"errors"
	"technical-test-backend/database"
	"technical-test-backend/models"

	"github.com/google/uuid"
)

type CatalogService struct{}

// AddToEtalaseInput - Input untuk seller menambahkan produk ke marketplace
type AddToEtalaseInput struct {
	ProductID    string  `json:"product_id" binding:"required"`    // UUID produk dari gudang pusat
	SellingPrice float64 `json:"selling_price" binding:"required,min=1"` // Harga jual seller
}

// MarketplaceItem - Struktur data untuk tampilan marketplace
type MarketplaceItem struct {
	ID            uuid.UUID `json:"seller_product_id"` // ID produk di etalase seller
	ProductName   string    `json:"product_name"`      // Nama produk
	Category      string    `json:"category"`          // Kategori produk
	SellerName    string    `json:"seller_name"`       // Nama toko seller
	Price         float64   `json:"price"`             // Harga jual
	StockTersedia int       `json:"stock_available"`   // Stok tersedia dari gudang pusat
}

// AddToEtalase - Seller menambahkan produk dari gudang pusat ke marketplace mereka
// Alur: Validasi produk exist -> Validasi harga jual >= harga modal -> Simpan ke SellerProduct
func (s *CatalogService) AddToEtalase(sellerID string, input AddToEtalaseInput) (models.SellerProduct, error) {
	sUUID, _ := uuid.Parse(sellerID)
	pUUID, _ := uuid.Parse(input.ProductID)

	// 1. Cek apakah produk master ada di gudang pusat
	var master models.Product
	if err := database.DB.First(&master, "id = ?", pUUID).Error; err != nil {
		return models.SellerProduct{}, errors.New("master product not found")
	}

	// 2. Validasi harga jual tidak boleh lebih rendah dari harga modal
	// Selling Price = Harga Modal + Markup Seller
	if input.SellingPrice < master.Price {
		return models.SellerProduct{}, errors.New("selling price lower than capital price")
	}

	// 3. Simpan produk ke etalase seller (marketplace)
	item := models.SellerProduct{
		SellerID:     sUUID,
		ProductID:    pUUID,
		SellingPrice: input.SellingPrice,
		IsActive:     true, // Default aktif
	}
	err := database.DB.Create(&item).Error
	return item, err
}

// GetMarketplaceItems - Get semua produk aktif di marketplace dengan filter
// Filter: search (nama produk), category, price range (min-max)
// Hanya menampilkan produk yang is_active = true
func (s *CatalogService) GetMarketplaceItems(search string, categoryID string, minPrice float64, maxPrice float64) ([]MarketplaceItem, error) {
	var items []models.SellerProduct
	
	// 1. Build query dengan base filter: hanya produk aktif
	query := database.DB.Preload("Product.ProductType").Preload("Seller").Where("is_active = ?", true)
	
	// 2. Filter pencarian berdasarkan nama produk (case insensitive)
	if search != "" {
		query = query.Joins("JOIN products ON seller_products.product_id = products.id").
			Where("products.name ILIKE ?", "%"+search+"%")
	}
	
	// 3. Filter berdasarkan kategori produk
	if categoryID != "" {
		query = query.Joins("JOIN products ON seller_products.product_id = products.id").
			Where("products.product_type_id = ?", categoryID)
	}
	
	// 4. Filter berdasarkan range harga
	if minPrice > 0 {
		query = query.Where("selling_price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("selling_price <= ?", maxPrice)
	}
	
	// 5. Execute query
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	var result []MarketplaceItem
	for _, item := range items {
		if item.Product.Stock > 0 {
			result = append(result, MarketplaceItem{
				ID:            item.ID,
				ProductName:   item.Product.Name,
				Category:      item.Product.ProductType.Name,
				SellerName:    item.Seller.Name,
				Price:         item.SellingPrice,
				StockTersedia: item.Product.Stock,
			})
		}
	}
	return result, nil
}

// Response structure for seller's product list
type SellerProductDetail struct {
	ID           string  `json:"id"`
	ProductName  string  `json:"product_name"`
	Category     string  `json:"category"`
	BasePrice    float64 `json:"base_price"`
	SellingPrice float64 `json:"selling_price"`
	ProfitMargin float64 `json:"profit_margin"`
	Stock        int     `json:"stock"`
	IsActive     bool    `json:"is_active"`
}

// Fungsi untuk melihat produk yang dijual oleh seller tertentu
func (s *CatalogService) GetSellerProducts(sellerID string) ([]SellerProductDetail, error) {
	var items []models.SellerProduct
	
	// Query Join untuk mendapatkan produk milik seller
	if err := database.DB.Preload("Product.ProductType").
		Where("seller_id = ?", sellerID).
		Find(&items).Error; err != nil {
		return nil, err
	}

	var result []SellerProductDetail
	for _, item := range items {
		profitMargin := ((item.SellingPrice - item.Product.Price) / item.Product.Price) * 100
		result = append(result, SellerProductDetail{
			ID:           item.ID.String(),
			ProductName:  item.Product.Name,
			Category:     item.Product.ProductType.Name,
			BasePrice:    item.Product.Price,
			SellingPrice: item.SellingPrice,
			ProfitMargin: profitMargin,
			Stock:        item.Product.Stock,
			IsActive:     item.IsActive,
		})
	}
	return result, nil
}

// Update Seller Product Price
type UpdateSellerProductInput struct {
	SellingPrice *float64 `json:"selling_price" binding:"omitempty,min=1"`
	IsActive     *bool    `json:"is_active"`
}

func (s *CatalogService) UpdateSellerProduct(sellerProductID string, sellerID string, input UpdateSellerProductInput) (models.SellerProduct, error) {
	var sellerProduct models.SellerProduct
	
	// Check if seller product exists and belongs to the seller
	if err := database.DB.Preload("Product").First(&sellerProduct, "id = ? AND seller_id = ?", sellerProductID, sellerID).Error; err != nil {
		return sellerProduct, errors.New("seller product not found or unauthorized")
	}

	updates := make(map[string]interface{})
	
	if input.SellingPrice != nil {
		// Validate price is not lower than base price
		if *input.SellingPrice < sellerProduct.Product.Price {
			return sellerProduct, errors.New("selling price cannot be lower than base price")
		}
		updates["selling_price"] = *input.SellingPrice
	}
	
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if err := database.DB.Model(&sellerProduct).Updates(updates).Error; err != nil {
		return sellerProduct, err
	}

	// Reload with associations
	database.DB.Preload("Product.ProductType").First(&sellerProduct, "id = ?", sellerProductID)
	return sellerProduct, nil
}

// Delete Seller Product (Soft delete - set is_active to false)
func (s *CatalogService) DeleteSellerProduct(sellerProductID string, sellerID string) error {
	var sellerProduct models.SellerProduct
	
	// Check if seller product exists and belongs to the seller
	if err := database.DB.First(&sellerProduct, "id = ? AND seller_id = ?", sellerProductID, sellerID).Error; err != nil {
		return errors.New("seller product not found or unauthorized")
	}

	// Soft delete by setting is_active to false
	return database.DB.Model(&sellerProduct).Update("is_active", false).Error
}