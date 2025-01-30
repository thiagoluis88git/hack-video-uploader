package model

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name  string
	CPF   string `gorm:"index;unique"`
	Email string `gorm:"unique"`
}
