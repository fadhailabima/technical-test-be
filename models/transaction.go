package models

import (
	"github.com/google/uuid"
)

const (
	StatusPending   = "PENDING"
	StatusCompleted = "COMPLETED"
	StatusCancelled = "CANCELLED"
)

type Transaction struct {
	Base
	UserID          uuid.UUID `gorm:"type:uuid;not null"` 
	SellerProductID uuid.UUID `gorm:"type:uuid;not null"` 
	Quantity        int       `gorm:"not null"`
	Status          string    `gorm:"type:varchar(20);default:'PENDING'"`
	TotalPrice   float64 `gorm:"type:decimal(15,2)"` 
	AdminFee     float64 `gorm:"type:decimal(15,2)"` 
	SellerProfit float64 `gorm:"type:decimal(15,2)"`

	User          User          `gorm:"foreignKey:UserID"`
	SellerProduct SellerProduct `gorm:"foreignKey:SellerProductID"`
}