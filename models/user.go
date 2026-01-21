package models

import (
	"github.com/google/uuid"
)

type User struct {
	Base
	Name     string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100);unique;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	RoleID   uuid.UUID `gorm:"type:uuid;not null"`
	Role     Role   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}