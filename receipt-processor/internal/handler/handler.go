package handler

import (
	"encoding/json"
	"net/http"
	"receipt-processor/internal/model"
	"regexp"
	"strconv"
	"time"
)

// IsValidPrice checks if a given price string is a valid price in terms of dollars and cents.
func IsValidPrice(price string) bool {
	// The price must start with one or more digits (\d+).
	// Optionally, it can have a decimal point followed by exactly two digits (\.\d{2})
	validPrice := regexp.MustCompile(`^\d+(\.\d{2})?$`)
	return validPrice.MatchString(price)
}

// ProcessReceipt handles HTTP requests for processing receipts. It validates the incoming receipt,
// computes the points associated with it, and stores it.
// Responds with the receipt ID.
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt model.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check for empty strings
	if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || len(receipt.Items) == 0 || receipt.Total == "" {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	// Check for invalid date
	t, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Check for future date
	if t.After(time.Now()) {
		http.Error(w, "Date cannot be in the future", http.StatusBadRequest)
		return
	}

	// Check for invalid time
	_, err = time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		http.Error(w, "Invalid time format", http.StatusBadRequest)
		return
	}

	for _, item := range receipt.Items {
		// Check for correct price format
		if !IsValidPrice(item.Price) {
			http.Error(w, "Invalid Price Format", http.StatusBadRequest)
			return
		}
		// Check for negative and zero prices in Items
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil || price == 0 {
			http.Error(w, "Zero Price error", http.StatusBadRequest)
			return
		}
	}

	// Check for correct price format
	if !IsValidPrice(receipt.Total) {
		http.Error(w, "Invalid Price Format", http.StatusBadRequest)
		return
	}

	// Check for negative or zero total price
	totalPrice, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil || totalPrice == 0 {
		http.Error(w, "Zero Price error", http.StatusBadRequest)
		return
	}

	receiptID := model.StoreReceipt(receipt)

	// Set response header and encode JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": receiptID})
}

// GetPoints handles HTTP requests for retrieving the points associated with a given receipt ID.
// Responds with the points or an error if the ID is not found.
func GetPoints(w http.ResponseWriter, r *http.Request, id string) {
	points, ok := model.GetPoints(id)

	if !ok {
		http.Error(w, "No receipt found for that id", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": points})
}
