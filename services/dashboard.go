package services

import (
	"technical-test-backend/database"
	"technical-test-backend/models"
	"time"
)

type DashboardService struct{}

// Customer Dashboard Response
type CustomerDashboard struct {
	TotalOrders     int64                     `json:"total_orders"`
	PendingOrders   int64                     `json:"pending_orders"`
	ConfirmedOrders int64                     `json:"confirmed_orders"`
	TotalSpent      float64                   `json:"total_spent"`
	RecentOrders    []CustomerRecentOrder     `json:"recent_orders"`
}

type CustomerRecentOrder struct {
	ID            string    `json:"id"`
	ProductName   string    `json:"product_name"`
	SellerName    string    `json:"seller_name"`
	Quantity      int       `json:"quantity"`
	TotalPrice    float64   `json:"total_price"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// Seller Dashboard Response
type SellerDashboard struct {
	ProductsInMarketplace int64             `json:"products_in_marketplace"`
	TotalSalesRevenue     float64           `json:"total_sales_revenue"`
	TotalTransactions     int64             `json:"total_transactions"`
	PendingOrders         int64             `json:"pending_orders"`
	ConfirmedOrders       int64             `json:"confirmed_orders"`
	TotalProfit           float64           `json:"total_profit"`
	ProfitMargin          float64           `json:"profit_margin_percentage"`
	TopProducts           []SellerTopProduct `json:"top_products"`
}

type SellerTopProduct struct {
	ProductName      string `json:"product_name"`
	TransactionCount int64  `json:"transaction_count"`
	TotalQuantity    int64  `json:"total_quantity"`
	TotalRevenue     float64 `json:"total_revenue"`
}

// Admin Dashboard Response
type AdminDashboard struct {
	TotalProducts      int64   `json:"total_products"`
	TotalProductTypes  int64   `json:"total_product_types"`
	TotalSellers       int64   `json:"total_sellers"`
	TotalCustomers     int64   `json:"total_customers"`
	TransactionsToday  int64   `json:"transactions_today"`
	PlatformIncome     float64 `json:"platform_income"`
}

// GetBuyerStats - Customer Dashboard
func (s *DashboardService) GetBuyerStats(userID string) CustomerDashboard {
	var stats CustomerDashboard
	
	// Total orders
	database.DB.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&stats.TotalOrders)
	
	// Pending orders
	database.DB.Model(&models.Transaction{}).Where("user_id = ? AND status = ?", userID, "PENDING").Count(&stats.PendingOrders)
	
	// Confirmed orders
	database.DB.Model(&models.Transaction{}).Where("user_id = ? AND status = ?", userID, models.StatusCompleted).Count(&stats.ConfirmedOrders)
	
	// Total spent (confirmed only)
	database.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND status = ?", userID, models.StatusCompleted).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&stats.TotalSpent)
	
	// Recent 3 orders with product details
	type OrderResult struct {
		TransactionID string
		ProductName   string
		SellerName    string
		Quantity      int
		TotalPrice    float64
		Status        string
		CreatedAt     time.Time
	}
	
	var orderResults []OrderResult
	database.DB.Table("transactions").
		Select("transactions.id as transaction_id, products.name as product_name, users.name as seller_name, transactions.quantity, transactions.total_price, transactions.status, transactions.created_at").
		Joins("JOIN seller_products ON transactions.seller_product_id = seller_products.id").
		Joins("JOIN products ON seller_products.product_id = products.id").
		Joins("JOIN users ON seller_products.seller_id = users.id").
		Where("transactions.user_id = ?", userID).
		Order("transactions.created_at DESC").
		Limit(3).
		Scan(&orderResults)
	
	for _, result := range orderResults {
		stats.RecentOrders = append(stats.RecentOrders, CustomerRecentOrder{
			ID:          result.TransactionID,
			ProductName: result.ProductName,
			SellerName:  result.SellerName,
			Quantity:    result.Quantity,
			TotalPrice:  result.TotalPrice,
			Status:      result.Status,
			CreatedAt:   result.CreatedAt,
		})
	}
	
	return stats
}

// GetSellerStats - Seller Dashboard
func (s *DashboardService) GetSellerStats(sellerID string) SellerDashboard {
	var stats SellerDashboard
	
	// Products in marketplace
	database.DB.Model(&models.SellerProduct{}).Where("seller_id = ?", sellerID).Count(&stats.ProductsInMarketplace)
	
	// Total sales revenue (confirmed transactions)
	database.DB.Table("transactions").
		Joins("JOIN seller_products ON transactions.seller_product_id = seller_products.id").
		Where("seller_products.seller_id = ? AND transactions.status = ?", sellerID, models.StatusCompleted).
		Select("COALESCE(SUM(transactions.total_price), 0)").
		Scan(&stats.TotalSalesRevenue)
	
	// Total transactions
	database.DB.Table("transactions").
		Joins("JOIN seller_products ON transactions.seller_product_id = seller_products.id").
		Where("seller_products.seller_id = ?", sellerID).
		Count(&stats.TotalTransactions)
	
	// Pending orders
	database.DB.Table("transactions").
		Joins("JOIN seller_products ON transactions.seller_product_id = seller_products.id").
		Where("seller_products.seller_id = ? AND transactions.status = ?", sellerID, "PENDING").
		Count(&stats.PendingOrders)
	
	// Confirmed orders
	database.DB.Table("transactions").
		Joins("JOIN seller_products ON transactions.seller_product_id = seller_products.id").
		Where("seller_products.seller_id = ? AND transactions.status = ?", sellerID, models.StatusCompleted).
		Count(&stats.ConfirmedOrders)
	
	// Total profit (seller_profit from confirmed transactions)
	database.DB.Table("transactions").
		Joins("JOIN seller_products ON transactions.seller_product_id = seller_products.id").
		Where("seller_products.seller_id = ? AND transactions.status = ?", sellerID, models.StatusCompleted).
		Select("COALESCE(SUM(transactions.seller_profit), 0)").
		Scan(&stats.TotalProfit)
	
	// Calculate profit margin percentage
	if stats.TotalSalesRevenue > 0 {
		stats.ProfitMargin = (stats.TotalProfit / stats.TotalSalesRevenue) * 100
	}
	
	// Top 3 products by transaction count
	type TopProductResult struct {
		ProductName      string
		TransactionCount int64
		TotalQuantity    int64
		TotalRevenue     float64
	}
	
	var topResults []TopProductResult
	database.DB.Table("transactions").
		Select("products.name as product_name, COUNT(transactions.id) as transaction_count, SUM(transactions.quantity) as total_quantity, SUM(transactions.total_price) as total_revenue").
		Joins("JOIN seller_products ON transactions.seller_product_id = seller_products.id").
		Joins("JOIN products ON seller_products.product_id = products.id").
		Where("seller_products.seller_id = ? AND transactions.status = ?", sellerID, models.StatusCompleted).
		Group("products.id, products.name").
		Order("transaction_count DESC").
		Limit(3).
		Scan(&topResults)
	
	for _, result := range topResults {
		stats.TopProducts = append(stats.TopProducts, SellerTopProduct{
			ProductName:      result.ProductName,
			TransactionCount: result.TransactionCount,
			TotalQuantity:    result.TotalQuantity,
			TotalRevenue:     result.TotalRevenue,
		})
	}
	
	return stats
}

// GetAdminStats - Admin Dashboard
func (s *DashboardService) GetAdminStats() AdminDashboard {
	var stats AdminDashboard
	
	// Total products
	database.DB.Model(&models.Product{}).Count(&stats.TotalProducts)
	
	// Total product types
	database.DB.Model(&models.ProductType{}).Count(&stats.TotalProductTypes)
	
	// Total sellers (count users with Seller role)
	database.DB.Table("users").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", "Seller").
		Count(&stats.TotalSellers)
	
	// Total customers (count users with Pelanggan role)
	database.DB.Table("users").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", "Pelanggan").
		Count(&stats.TotalCustomers)
	
	// Transactions today
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.Transaction{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.TransactionsToday)
	
	// Platform income (admin_fee from confirmed transactions)
	database.DB.Model(&models.Transaction{}).
		Where("status = ?", models.StatusCompleted).
		Select("COALESCE(SUM(admin_fee), 0)").
		Scan(&stats.PlatformIncome)
	
	return stats
}