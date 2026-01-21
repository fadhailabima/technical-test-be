package services

import (
	"errors"
	"technical-test-backend/database"
	"technical-test-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type TransactionService struct{}

type CreateOrderInput struct {
	SellerProductID string `json:"seller_product_id" binding:"required"`
	Quantity        int    `json:"quantity" binding:"required,min=1"`
	UserID          string // Dari Token
}

func (s *TransactionService) CreateOrder(input CreateOrderInput) (models.Transaction, error) {
	sellerProductUUID, _ := uuid.Parse(input.SellerProductID)
	userUUID, _ := uuid.Parse(input.UserID)

	// Ambil Data SellerProduct + Data Product Asli (Admin)
	var item models.SellerProduct
	if err := database.DB.Preload("Product").First(&item, "id = ?", sellerProductUUID).Error; err != nil {
		return models.Transaction{}, errors.New("barang tidak ditemukan")
	}

	// Validasi Stok Tersedia
	if item.Product.Stock < input.Quantity {
		return models.Transaction{}, errors.New("stok tidak mencukupi")
	}

	// Validasi SellerProduct Aktif
	if !item.IsActive {
		return models.Transaction{}, errors.New("produk tidak aktif")
	}

	// Hitung Kalkulasi Keuangan
	qty := float64(input.Quantity)
	
	// Uang Masuk dari Pembeli (Harga Seller * Qty)
	grandTotal := item.SellingPrice * qty 
	
	// Jatah Admin (Harga Modal * Qty)
	totalAdminFee := item.Product.Price * qty 
	
	// Jatah Seller (Sisa uang)
	totalSellerProfit := grandTotal - totalAdminFee

	// Simpan Transaksi dengan Detail Keuangan
	transaction := models.Transaction{
		UserID:          userUUID,
		SellerProductID: sellerProductUUID,
		Quantity:        input.Quantity,
		Status:          models.StatusPending,
		
		// Simpan Snapshot Keuangan
		TotalPrice:   grandTotal,
		AdminFee:     totalAdminFee,
		SellerProfit: totalSellerProfit,
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}
// GET Customer Transactions - List semua transaksi pembelian customer
type CustomerTransactionDetail struct {
	ID             string  `json:"id"`
	ProductName    string  `json:"product_name"`
	SellerName     string  `json:"seller_name"`
	SellerEmail    string  `json:"seller_email"`
	Quantity       int     `json:"quantity"`
	TotalPrice     float64 `json:"total_price"`
	AdminFee       float64 `json:"admin_fee"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
}

func (s *TransactionService) GetCustomerTransactions(customerID string) ([]CustomerTransactionDetail, error) {
	customerUUID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, errors.New("invalid customer ID")
	}

	var results []struct {
		TransactionID string
		ProductName   string
		SellerName    string
		SellerEmail   string
		Quantity      int
		TotalPrice    float64
		AdminFee      float64
		Status        string
		CreatedAt     string
	}

	err = database.DB.Table("transactions").
		Select(`
			transactions.id as transaction_id,
			products.name as product_name,
			users.name as seller_name,
			users.email as seller_email,
			transactions.quantity,
			transactions.total_price,
			transactions.admin_fee,
			transactions.status,
			transactions.created_at
		`).
		Joins("JOIN seller_products ON transactions.seller_product_id = seller_products.id").
		Joins("JOIN products ON seller_products.product_id = products.id").
		Joins("JOIN users ON seller_products.seller_id = users.id").
		Where("transactions.user_id = ?", customerUUID).
		Order("transactions.created_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	var transactions []CustomerTransactionDetail
	for _, r := range results {
		transactions = append(transactions, CustomerTransactionDetail{
			ID:          r.TransactionID,
			ProductName: r.ProductName,
			SellerName:  r.SellerName,
			SellerEmail: r.SellerEmail,
			Quantity:    r.Quantity,
			TotalPrice:  r.TotalPrice,
			AdminFee:    r.AdminFee,
			Status:      r.Status,
			CreatedAt:   r.CreatedAt,
		})
	}

	return transactions, nil
}
// SELLER CONFIRM (Potong Stok Admin)
func (s *TransactionService) ConfirmOrder(transactionID string, sellerID string) error {
	txUUID, _ := uuid.Parse(transactionID)

	txDB := database.DB.Begin()

	// Cek Transaksi milik Seller ini
	var transaction models.Transaction
	if err := txDB.Preload("SellerProduct").
		Joins("JOIN seller_products ON seller_products.id = transactions.seller_product_id").
		Where("transactions.id = ? AND seller_products.seller_id = ?", txUUID, sellerID).
		First(&transaction).Error; err != nil {
		txDB.Rollback()
		return errors.New("transaksi tidak ditemukan atau akses ditolak")
	}

	if transaction.Status != models.StatusPending {
		txDB.Rollback()
		return errors.New("transaksi sudah selesai/batal")
	}

	// Kunci & Cek Stok Admin (Locking)
	var masterProduct models.Product
	if err := txDB.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&masterProduct, "id = ?", transaction.SellerProduct.ProductID).Error; err != nil {
		txDB.Rollback()
		return errors.New("produk master hilang")
	}

	if masterProduct.Stock < transaction.Quantity {
		// Auto Cancel jika stok admin habis
		transaction.Status = models.StatusCancelled
		txDB.Save(&transaction)
		txDB.Commit()
		return errors.New("stok gudang pusat habis")
	}

	// Kurangi Stok & Selesaikan
	masterProduct.Stock -= transaction.Quantity
	transaction.Status = models.StatusCompleted

	if err := txDB.Save(&masterProduct).Error; err != nil {
		txDB.Rollback(); return err
	}
	if err := txDB.Save(&transaction).Error; err != nil {
		txDB.Rollback(); return err
	}
	
	txDB.Commit()
	return nil
}

// GET Seller Transactions - List semua transaksi dari produk seller
type SellerTransactionDetail struct {
	ID              string  `json:"id"`
	ProductName     string  `json:"product_name"`
	BuyerName       string  `json:"buyer_name"`
	BuyerEmail      string  `json:"buyer_email"`
	Quantity        int     `json:"quantity"`
	TotalPrice      float64 `json:"total_price"`
	SellerProfit    float64 `json:"seller_profit"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}

func (s *TransactionService) GetSellerTransactions(sellerID string) ([]SellerTransactionDetail, error) {
	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return nil, errors.New("invalid seller ID")
	}

	type QueryResult struct {
		TransactionID string
		ProductName   string
		BuyerName     string
		BuyerEmail    string
		Quantity      int
		TotalPrice    float64
		SellerProfit  float64
		Status        string
		CreatedAt     string
	}

	var results []QueryResult
	err = database.DB.Table("transactions").
		Select(`
			transactions.id as transaction_id,
			products.name as product_name,
			users.name as buyer_name,
			users.email as buyer_email,
			transactions.quantity,
			transactions.total_price,
			transactions.seller_profit,
			transactions.status,
			transactions.created_at
		`).
		Joins("JOIN seller_products ON transactions.seller_product_id = seller_products.id").
		Joins("JOIN products ON seller_products.product_id = products.id").
		Joins("JOIN users ON transactions.user_id = users.id").
		Where("seller_products.seller_id = ?", sellerUUID).
		Order("transactions.created_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	var transactions []SellerTransactionDetail
	for _, r := range results {
		transactions = append(transactions, SellerTransactionDetail{
			ID:           r.TransactionID,
			ProductName:  r.ProductName,
			BuyerName:    r.BuyerName,
			BuyerEmail:   r.BuyerEmail,
			Quantity:     r.Quantity,
			TotalPrice:   r.TotalPrice,
			SellerProfit: r.SellerProfit,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
		})
	}

	return transactions, nil
}

// GetTransactionDetail - Get single transaction by ID
type TransactionDetail struct {
	ID             string  `json:"id"`
	ProductName    string  `json:"product_name"`
	BuyerName      string  `json:"buyer_name"`
	BuyerEmail     string  `json:"buyer_email"`
	SellerName     string  `json:"seller_name"`
	SellerEmail    string  `json:"seller_email"`
	Quantity       int     `json:"quantity"`
	TotalPrice     float64 `json:"total_price"`
	AdminFee       float64 `json:"admin_fee"`
	SellerProfit   float64 `json:"seller_profit"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
}

func (s *TransactionService) GetTransactionDetail(transactionID string) (TransactionDetail, error) {
	txUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return TransactionDetail{}, errors.New("invalid transaction ID")
	}

	var result struct {
		TransactionID string
		ProductName   string
		BuyerName     string
		BuyerEmail    string
		SellerName    string
		SellerEmail   string
		Quantity      int
		TotalPrice    float64
		AdminFee      float64
		SellerProfit  float64
		Status        string
		CreatedAt     string
	}

	err = database.DB.Table("transactions").
		Select(`
			transactions.id as transaction_id,
			products.name as product_name,
			buyer.name as buyer_name,
			buyer.email as buyer_email,
			seller.name as seller_name,
			seller.email as seller_email,
			transactions.quantity,
			transactions.total_price,
			transactions.admin_fee,
			transactions.seller_profit,
			transactions.status,
			transactions.created_at
		`).
		Joins("JOIN seller_products ON transactions.seller_product_id = seller_products.id").
		Joins("JOIN products ON seller_products.product_id = products.id").
		Joins("JOIN users as buyer ON transactions.user_id = buyer.id").
		Joins("JOIN users as seller ON seller_products.seller_id = seller.id").
		Where("transactions.id = ?", txUUID).
		Scan(&result).Error

	if err != nil {
		return TransactionDetail{}, err
	}

	return TransactionDetail{
		ID:           result.TransactionID,
		ProductName:  result.ProductName,
		BuyerName:    result.BuyerName,
		BuyerEmail:   result.BuyerEmail,
		SellerName:   result.SellerName,
		SellerEmail:  result.SellerEmail,
		Quantity:     result.Quantity,
		TotalPrice:   result.TotalPrice,
		AdminFee:     result.AdminFee,
		SellerProfit: result.SellerProfit,
		Status:       result.Status,
		CreatedAt:    result.CreatedAt,
	}, nil
}

// CancelTransaction - Cancel transaction by customer
func (s *TransactionService) CancelTransaction(transactionID string, userID string) error {
	txUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return errors.New("invalid transaction ID")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Find transaction
	var transaction models.Transaction
	if err := database.DB.Where("id = ? AND user_id = ?", txUUID, userUUID).First(&transaction).Error; err != nil {
		return errors.New("transaction not found or unauthorized")
	}

	// Only allow cancellation for PENDING transactions
	if transaction.Status != models.StatusPending {
		return errors.New("only pending transactions can be cancelled")
	}

	// Update status to CANCELLED
	transaction.Status = models.StatusCancelled
	if err := database.DB.Save(&transaction).Error; err != nil {
		return err
	}

	return nil
}