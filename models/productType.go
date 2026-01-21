package models


type ProductType struct {
	Base
	Name string `gorm:"type:varchar(50);not null;unique"` 
}