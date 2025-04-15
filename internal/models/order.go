package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
	Side     string  `json:"side"`   // "buy" or "sell"
	Status   string  `json:"status"` // "open", "closed", "canceled"
}
