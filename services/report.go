package services

import (
	"technical-test-backend/database"
	"time"
)

type ReportService struct{}

// Sales Report
type SalesReportItem struct {
	Date          string  `json:"date"`
	TotalOrders   int     `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
	AdminIncome   float64 `json:"admin_income"`
	SellerIncome  float64 `json:"seller_income"`
}

func (s *ReportService) GetSalesReport(startDate, endDate string) ([]SalesReportItem, error) {
	var results []struct {
		Date         string
		TotalOrders  int
		TotalRevenue float64
		AdminIncome  float64
		SellerIncome float64
	}

	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as total_orders,
			SUM(total_price) as total_revenue,
			SUM(admin_fee) as admin_income,
			SUM(seller_profit) as seller_income
		FROM transactions
		WHERE status = 'COMPLETED'
	`

	if startDate != "" && endDate != "" {
		query += " AND DATE(created_at) BETWEEN ? AND ?"
		err := database.DB.Raw(query+" GROUP BY DATE(created_at) ORDER BY date DESC", startDate, endDate).Scan(&results).Error
		if err != nil {
			return nil, err
		}
	} else {
		// Default last 30 days
		endDateParsed := time.Now()
		startDateParsed := endDateParsed.AddDate(0, 0, -30)
		query += " AND DATE(created_at) BETWEEN ? AND ?"
		err := database.DB.Raw(query+" GROUP BY DATE(created_at) ORDER BY date DESC", startDateParsed.Format("2006-01-02"), endDateParsed.Format("2006-01-02")).Scan(&results).Error
		if err != nil {
			return nil, err
		}
	}

	var report []SalesReportItem
	for _, r := range results {
		report = append(report, SalesReportItem{
			Date:          r.Date,
			TotalOrders:   r.TotalOrders,
			TotalRevenue:  r.TotalRevenue,
			AdminIncome:   r.AdminIncome,
			SellerIncome:  r.SellerIncome,
		})
	}

	return report, nil
}

// Top Products Report
type TopProductItem struct {
	ProductName       string  `json:"product_name"`
	Category          string  `json:"category"`
	TotalSold         int     `json:"total_sold"`
	TotalRevenue      float64 `json:"total_revenue"`
	TotalTransactions int     `json:"total_transactions"`
}

func (s *ReportService) GetTopProducts(limit int) ([]TopProductItem, error) {
	if limit <= 0 {
		limit = 10
	}

	var results []struct {
		ProductName       string
		Category          string
		TotalSold         int
		TotalRevenue      float64
		TotalTransactions int
	}

	query := `
		SELECT 
			products.name as product_name,
			product_types.name as category,
			SUM(transactions.quantity) as total_sold,
			SUM(transactions.total_price) as total_revenue,
			COUNT(transactions.id) as total_transactions
		FROM transactions
		JOIN seller_products ON transactions.seller_product_id = seller_products.id
		JOIN products ON seller_products.product_id = products.id
		JOIN product_types ON products.product_type_id = product_types.id
		WHERE transactions.status = 'COMPLETED'
		GROUP BY products.id, products.name, product_types.name
		ORDER BY total_sold DESC
		LIMIT ?
	`

	err := database.DB.Raw(query, limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var report []TopProductItem
	for _, r := range results {
		report = append(report, TopProductItem{
			ProductName:       r.ProductName,
			Category:          r.Category,
			TotalSold:         r.TotalSold,
			TotalRevenue:      r.TotalRevenue,
			TotalTransactions: r.TotalTransactions,
		})
	}

	return report, nil
}

// Top Sellers Report
type TopSellerItem struct {
	SellerName        string  `json:"seller_name"`
	SellerEmail       string  `json:"seller_email"`
	TotalProducts     int     `json:"total_products"`
	TotalSales        float64 `json:"total_sales"`
	TotalProfit       float64 `json:"total_profit"`
	TotalTransactions int     `json:"total_transactions"`
}

func (s *ReportService) GetTopSellers(limit int) ([]TopSellerItem, error) {
	if limit <= 0 {
		limit = 10
	}

	var results []struct {
		SellerName        string
		SellerEmail       string
		TotalProducts     int
		TotalSales        float64
		TotalProfit       float64
		TotalTransactions int
	}

	query := `
		SELECT 
			users.name as seller_name,
			users.email as seller_email,
			COUNT(DISTINCT seller_products.product_id) as total_products,
			SUM(transactions.total_price) as total_sales,
			SUM(transactions.seller_profit) as total_profit,
			COUNT(transactions.id) as total_transactions
		FROM users
		JOIN seller_products ON users.id = seller_products.seller_id
		JOIN transactions ON seller_products.id = transactions.seller_product_id
		WHERE transactions.status = 'COMPLETED'
		GROUP BY users.id, users.name, users.email
		ORDER BY total_sales DESC
		LIMIT ?
	`

	err := database.DB.Raw(query, limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var report []TopSellerItem
	for _, r := range results {
		report = append(report, TopSellerItem{
			SellerName:        r.SellerName,
			SellerEmail:       r.SellerEmail,
			TotalProducts:     r.TotalProducts,
			TotalSales:        r.TotalSales,
			TotalProfit:       r.TotalProfit,
			TotalTransactions: r.TotalTransactions,
		})
	}

	return report, nil
}
