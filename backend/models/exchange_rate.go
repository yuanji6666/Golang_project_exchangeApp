package models

import (
	"time"
)

type ExchangeRate struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	FromCurrency string    `gorm:"column:from_currency" json:"from_currency" binding:"required"`
	ToCurrency   string    `gorm:"column:to_currency" json:"to_currency" binding:"required"`
	Rate         float64   `json:"rate" binding:"required"`
	Date         time.Time `json:"date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
