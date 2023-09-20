package model

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// In-memory storage
var receipts = make(map[string]Receipt)
var receiptPoints = make(map[string]int)
var mu sync.Mutex // Mutex for locking access to the maps

// StoreReceipt saves a receipt and returns a generated ID
func StoreReceipt(receipt Receipt) string {
	mu.Lock()
	defer mu.Unlock()

	// Combine current time and a random number for the ID to avoid collisions
	id := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000000))

	receipts[id] = receipt
	points := TallyPoints(receipt)
	receiptPoints[id] = points
	return id
}

func TallyPoints(receipt Receipt) int {

	var points int = 0

	// Add points for all alphanumeric characters
	is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	for _, char := range receipt.Retailer {
		if is_alphanumeric(string(char)) {
			points++
		}
	}

	// Add points for multiples of 0.25 and for no cent dollar amounts
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil {
		totalCents := int(totalFloat * 100)
		if totalCents%100 == 0 {
			points += 50
		}

		if totalCents%25 == 0 {
			points += 25
		}

	}

	var itemCount int = 0
	for _, item := range receipt.Items {

		// Points for itemCounts that are at 2
		itemCount++
		if itemCount == 2 {
			itemCount = 0
			points += 5
		}

		// Points for item descriptions being multiples of 3
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				additionalPoints := math.Ceil(price * 0.2)
				points += int(additionalPoints)
			}
		}
	}

	// Points if the day in the purchase date is odd
	t, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if t.Day()%2 == 1 {
		points += 6
	}

	// Points if the time of purchase is after 2:00pm and before 4:00pm
	t, _ = time.Parse("15:04", receipt.PurchaseTime)
	minutesPastMidnight := t.Hour()*60 + t.Minute() // t.hour only gives hours, needs minutes too
	if minutesPastMidnight > (14*60) && minutesPastMidnight < (16*60) {
		points += 10
	}

	return points
}

// GetPoints retrieves points for a receipt ID
func GetPoints(id string) (int, bool) {
	mu.Lock()
	defer mu.Unlock()

	points, ok := receiptPoints[id]
	return points, ok
}
