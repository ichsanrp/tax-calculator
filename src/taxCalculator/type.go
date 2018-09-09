package taxCalculator

import (
	"time"
)

const (
	FOOD_TAX          = 1
	TOBACCO_TAX       = 2
	ENTERTAINMENT_TAX = 3
)

type Item struct {
	ID        int     `json:"id"`
	Weight    int     `json:"weight"`
	TaxCodeID int     `json:"tax_code"`
	SessionID int     `json:"session_id"`
	Price     float32 `json:"price"`
	Tax       float32 `json:"tax"`
	Name      string  `json:"name"`
}

type Session struct {
	ID         int       `json:"id"`
	CreateTime time.Time `json:"create_time"`
}

type DefaultResponse struct {
	Header *DefaultResponseHeader `json:"header"`
	Data   interface{}            `json:"data"`
}

type DefaultResponseHeader struct {
	ProcessTime float32 `json:"process_time"`
	Error       string  `json:"error"`
	StatusCode  int     `json:"status_code"`
}
