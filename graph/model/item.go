package model

type DiscountType string

const (
	BogoDiscountType       DiscountType = "BOGO"
	PercentageDiscountType DiscountType = "PERCENTAGE"
	NoneDiscountType       DiscountType = "NONE"
	DollarsOffDiscountType DiscountType = "DOLLARS_OFF"
	FixedPriceDiscountType DiscountType = "FIXED_PRICE"
)

func (dt DiscountType) String() string {
	return string(dt)
}

type Item struct {
	Code         string       `json:"code"`
	Price        float64      `json:"price"`
	Weight       *float64     `json:"weight,omitempty"`
	Savings      float64      `json:"savings"`
	Quantity     int          `json:"quantity"`
	ExtraInfo    *string      `json:"extra_info,omitempty"`
	Restricted   bool         `json:"restricted"`
	Description  string       `json:"description"`
	ProductKey   string       `json:"product_key"`
	TotalPrice   float64      `json:"total_price"`
	PriceWeight  *float64     `json:"price_weight,omitempty"`
	TaxRate      *float64     `json:"tax_rate,omitempty"`
	DiscountType DiscountType `json:"discount_type"`
	Discount     *string      `json:"discount,omitempty"`
	InternalKey  string       `json:"internal_key"`
	BogoID       *int         `json:"bogo_id,omitempty"`
	Category     *string      `json:"category,omitempty"`
}
