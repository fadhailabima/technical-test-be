package models

import (
	"github.com/google/uuid" 
)

type Product struct {
	Base
	Name          string      `gorm:"type:varchar(100);not null"`
	Stock         int         `gorm:"not null;check:stock >= 0"`
	Price         float64     `gorm:"type:decimal(10,2);not null"`
	ProductTypeID uuid.UUID     `gorm:"type:uuid;not null"`
	ProductType   ProductType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	
	CreatedAt     int64       `gorm:"autoCreateTime"`
}