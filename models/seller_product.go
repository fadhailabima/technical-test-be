package models

import (
	"github.com/google/uuid"
)

type SellerProduct struct {
	Base
	SellerID  uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"` 
	SellingPrice float64 `gorm:"type:decimal(15,2);not null"` 
	
	IsActive  bool `gorm:"default:true"`

	Seller  User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}