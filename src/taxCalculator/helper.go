package taxCalculator

import (
	"net/http"
	"strconv"
)

// calculateTax is function to calculat tax based on this rule
// food tax 10% of price
// Tobacco tax 20 + (2% * Price)
// entertainment tax is 1% from price - 100
func calculateTax(price float32, taxCode int) (tax float32) {
	switch taxCode {
	case FOOD_TAX:
		tax = float32(price) * 0.1
	case TOBACCO_TAX:
		tax = 20 + (0.02 * price)
	case ENTERTAINMENT_TAX:
		if price >= 100 {
			tax = 0.01 * (price - 100)
		}
	}
	return
}

// parseItemFromRequest used for translating http request parameters to Item object
func parseItemFromRequest(r *http.Request) (item *Item) {
	item = &Item{
		Name: r.FormValue("name"),
	}

	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 32)
	item.ID = int(id)

	price, _ := strconv.ParseFloat(r.FormValue("price"), 32)
	item.Price = float32(price)

	taxCode, _ := strconv.ParseInt(r.FormValue("tax_code"), 10, 32)
	item.TaxCodeID = int(taxCode)

	tax, _ := strconv.ParseFloat(r.FormValue("tax"), 32)
	item.Tax = float32(tax)

	sessionID, _ := strconv.ParseInt(r.FormValue("session_id"), 10, 32)
	item.SessionID = int(sessionID)

	return
}
