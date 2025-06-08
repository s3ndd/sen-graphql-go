package model

import (
	"time"
)

type Session struct {
	ID         string  `json:"id"`
	RetailerID *string `json:"retailer_id,omitempty"`
	CartID     string  `json:"cart_id"`
	Cart       *Cart   `json:"cart,omitempty"`
	UserID     *string `json:"user_id"`
	SiteID     string  `json:"site_id"`
	Status     string  `json:"status"`
	IntegrationType
	ExternalToken *string   `json:"external_token"`
	DealBarcode   *string   `json:"deal_barcode"`
	Items         []Item    `json:"items"`
	ItemsPrePos   []Item    `json:"items_pre_pos"`
	Total         float64   `json:"total"`
	TotalSavings  float64   `json:"total_savings"`
	TotalTax      []Tax     `json:"tax"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
}

type Tax struct {
	Rate   float64 `json:"rate"`
	Amount float64 `json:"amount"`
}
