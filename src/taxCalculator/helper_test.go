package taxCalculator

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func Test_calculateTax(t *testing.T) {
	type args struct {
		price   float32
		taxCode int
	}
	tests := []struct {
		name    string
		args    args
		wantTax float32
	}{
		{
			name: "pizza",
			args: args{
				price:   200,
				taxCode: 1,
			},
			wantTax: 20,
		},
		{
			name: "malboro",
			args: args{
				price:   100,
				taxCode: 2,
			},
			wantTax: 22,
		},
		{
			name: "movie ticket",
			args: args{
				price:   100,
				taxCode: 3,
			},
			wantTax: 0,
		},
		{
			name: "club ticket",
			args: args{
				price:   200,
				taxCode: 3,
			},
			wantTax: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTax := calculateTax(tt.args.price, tt.args.taxCode); gotTax != tt.wantTax {
				t.Errorf("calculateTax() = %v, want %v", gotTax, tt.wantTax)
			}
		})
	}
}

func Test_parseItemFromRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name     string
		args     args
		wantItem *Item
	}{
		{
			name: "pizza",
			args: args{
				r: func() (req *http.Request) {
					body := []byte("id=2&name=pizza&price=200&tax=20&tax_code=1&session_id=1")
					req, _ = http.NewRequest("POST", "http://localhost/item", bytes.NewBuffer(body))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					return
				}(),
			},
			wantItem: &Item{
				Name:      "pizza",
				Price:     200,
				TaxCodeID: 1,
				SessionID: 1,
				Tax:       20,
				ID:        2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotItem := parseItemFromRequest(tt.args.r); !reflect.DeepEqual(gotItem, tt.wantItem) {
				t.Errorf("parseItemFromRequest() = %v, want %v", gotItem, tt.wantItem)
			}
		})
	}
}
