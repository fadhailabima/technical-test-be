package models

type Role struct {
	Base
	Name string `gorm:"type:varchar(50);unique;not null"`
}