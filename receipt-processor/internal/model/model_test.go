package model_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-processor/internal/handler"
	"receipt-processor/internal/model"
	"testing"
)

// Testing function for Tallying points
func TestTallyPoints(t *testing.T) {
	testCases := []struct {
		retailer       string
		purchaseDate   string
		purchaseTime   string
		total          string
		items          []model.Item
		expectedPoints int
	}{
		// tests for alphanumerics
		{"ABC", "2023-09-19", "15:00", "10.00", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 94},
		{"123", "2023-09-19", "15:00", "10.00", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 94},
		{"]'..,'.'/][][.].,", "2023-09-19", "15:00", "10.00", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 91},
		{"123,ABC,real,=-';][;]", "2023-09-19", "15:00", "10.00", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 101},
		{"REAL,123,talk,.", "2023-09-18", "12:00", "0.01", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 11},
		// tests for time
		{"XYZ", "2023-09-19", "13:59", "10.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 34},
		{"XYZ", "2023-09-19", "14:00", "10.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 34},
		{"XYZ", "2023-09-19", "14:01", "10.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 44},
		{"XYZ", "2023-09-19", "15:00", "10.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 44},
		{"XYZ", "2023-09-19", "16:01", "10.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 34},
		{"XYZ", "2023-09-19", "16:00", "10.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 34},
		{"XYZ", "2023-09-19", "15:59", "10.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 44},
		// tests on multiples of 0.25 and round dollar values
		{"XYZ", "2023-09-19", "15:59", "0.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 44},
		{"XYZ", "2023-09-19", "15:59", "10.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 44},
		{"XYZ", "2023-09-19", "15:59", "10.26", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 19},
		{"XYZ", "2023-09-19", "15:59", "0.33", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 19},
		{"XYZ", "2023-09-19", "15:59", "1.00", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 94},
		{"XYZ", "2023-09-19", "15:59", "10.00", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 94},
		{"XYZ", "2023-09-19", "15:59", "1", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 94},
		// tests on purchase date
		{"XYZ", "2023-09-19", "15:59", "0.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 44},
		{"XYZ", "2023-09-18", "15:59", "0.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 38},
		{"XYZ", "2023-09-30", "15:59", "0.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 38},
		{"XYZ", "2023-09-01", "15:59", "0.25", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 44},
		// tests on item count
		{"[]", "2023-09-18", "12:00", "0.01", []model.Item{{Price: "0.01", ShortDescription: "abccd"}}, 0},
		{"[]", "2023-09-18", "12:00", "0.01", []model.Item{{Price: "0.01", ShortDescription: "abccd"}, {Price: "0.01", ShortDescription: "abccd"}}, 5},
		{"[]", "2023-09-18", "12:00", "0.01", []model.Item{{Price: "0.01", ShortDescription: "abccd"}, {Price: "0.01", ShortDescription: "abccd"}, {Price: "0.01", ShortDescription: "abccd"}}, 5},
		{"[]", "2023-09-18", "12:00", "0.01", []model.Item{{Price: "0.01", ShortDescription: "abccd"}, {Price: "0.01", ShortDescription: "abccd"}, {Price: "0.01", ShortDescription: "abccd"}, {Price: "0.01", ShortDescription: "abccd"}}, 10},

		//tests on item descriptions that are multiples of 3
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "1.20", ShortDescription: "abccd"}}, 0},
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "1.20", ShortDescription: "abccd"}, {Price: "1.20", ShortDescription: "abccd"}}, 5},
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "1.20", ShortDescription: "abccd"}, {Price: "1.20", ShortDescription: "abccd"}, {Price: "1.20", ShortDescription: "abccd"}, {Price: "1.20", ShortDescription: "abccd"}}, 10},
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "1.20", ShortDescription: "abc"}, {Price: "1.20", ShortDescription: "abc"}}, 7},
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "1.20", ShortDescription: "abc"}, {Price: "1.20", ShortDescription: "abc"}, {Price: "1.20", ShortDescription: "abc"}, {Price: "1.20", ShortDescription: "abc"}}, 14},
		// trimmed should mean removing trailing and leading whitespaces, meaning whitespaces in the middle should count
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "1.00", ShortDescription: "abc def"}}, 0},
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "1.00", ShortDescription: " abcdef "}}, 1},
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "5.00", ShortDescription: "abcdefghi"}}, 1},
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "4.99", ShortDescription: "aaa    aa"}}, 1},
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "0", ShortDescription: "avc"}}, 0},
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "10", ShortDescription: "abc"}, {Price: "100", ShortDescription: "abc"}, {Price: "0.01", ShortDescription: "abc"}, {Price: "1.20", ShortDescription: "abc"}}, 34},
		{"]", "2023-09-18", "14:00", "0.01", []model.Item{{Price: "100", ShortDescription: "a"}, {Price: "100", ShortDescription: "abc"}, {Price: "1.20", ShortDescription: "abcdef"}, {Price: "1.20", ShortDescription: "abcdefg"}}, 31},
	}

	for _, tc := range testCases {
		receipt := model.Receipt{
			Retailer:     tc.retailer,
			PurchaseDate: tc.purchaseDate,
			PurchaseTime: tc.purchaseTime,
			Total:        tc.total,
			Items:        tc.items,
		}
		id := model.StoreReceipt(receipt)
		points, ok := model.GetPoints(id)
		if !ok {
			t.Errorf("No points stored for receipt ID %s", id)
			continue
		}

		if points != tc.expectedPoints {
			t.Errorf("For receipt %+v, expected %d points, got %d points", receipt, tc.expectedPoints, points)
		}
	}
}

// test function for failure conditions in the process Endpoint
func TestProcessReceipt_Failures(t *testing.T) {
	testCases := []struct {
		name       string
		input      model.Receipt
		inputRaw   []byte
		httpStatus int
		errorMsg   string
	}{
		{
			name:       "Missing or invalid fields",
			input:      model.Receipt{},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Missing or invalid fields\n",
		},
		{
			name: "Missing or invalid fields",
			input: model.Receipt{
				Retailer:     "", //empty
				PurchaseDate: "2023-13-01",
				PurchaseTime: "15:00",
				Items:        []model.Item{{"item1", "2.50"}},
				Total:        "5.00",
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Missing or invalid fields\n",
		},
		{
			name: "Missing or invalid fields",
			input: model.Receipt{
				Retailer:     "dfsaf",
				PurchaseDate: "", // empty
				PurchaseTime: "15:00",
				Items:        []model.Item{{"item1", "2.50"}},
				Total:        "5.00",
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Missing or invalid fields\n",
		},
		{
			name: "Missing or invalid fields",
			input: model.Receipt{
				Retailer:     "dfsaf",
				PurchaseDate: "2023-13-01",
				PurchaseTime: "", // empty
				Items:        []model.Item{{"item1", "2.50"}},
				Total:        "5.00",
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Missing or invalid fields\n",
		},
		{
			name: "Missing or invalid fields",
			input: model.Receipt{
				Retailer:     "dfsaf",
				PurchaseDate: "2023-13-01",
				PurchaseTime: "15:00",
				Items:        []model.Item{}, // empty
				Total:        "5.00",
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Missing or invalid fields\n",
		},
		{
			name: "Missing or invalid fields",
			input: model.Receipt{
				Retailer:     "dfsaf",
				PurchaseDate: "2023-13-01",
				PurchaseTime: "15:00",
				Items:        []model.Item{{"item1", "2.50"}},
				Total:        "", // empty
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Missing or invalid fields\n",
		},
		{
			name:       "Invalid request payload",
			inputRaw:   []byte(`{"invalid json`),
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Invalid request payload\n",
		},
		{
			name: "Invalid date format",
			input: model.Receipt{
				Retailer:     "Walmart",
				PurchaseDate: "2023-13-01", // invalid date
				PurchaseTime: "15:00",
				Items:        []model.Item{{"item1", "2.50"}},
				Total:        "5.00",
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Invalid date format\n",
		},
		{
			name: "Date cannot be in the future",
			input: model.Receipt{
				Retailer:     "Walmart",
				PurchaseDate: "2023-11-10", // Future date
				PurchaseTime: "15:00",
				Items:        []model.Item{{"item1", "2.50"}},
				Total:        "5.00",
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Date cannot be in the future\n",
		},

		{
			name: "Invalid time format",
			input: model.Receipt{
				Retailer:     "Walmart",
				PurchaseDate: "2023-08-10",
				PurchaseTime: "25:00", // invalid time
				Items:        []model.Item{{"item1", "2.50"}},
				Total:        "5.00",
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Invalid time format\n",
		},
		{
			name: "Invalid Item Price negative",
			input: model.Receipt{
				Retailer:     "Walmart",
				PurchaseDate: "2023-08-10",
				PurchaseTime: "14:00",
				Items:        []model.Item{{"item1", "-2.50"}}, // negative item price
				Total:        "10.00",
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Invalid Price Format\n",
		},
		{
			name: "Invalid Item Price Zero",
			input: model.Receipt{
				Retailer:     "Walmart",
				PurchaseDate: "2023-08-10",
				PurchaseTime: "14:00",
				Items:        []model.Item{{"item1", "0"}}, // 0 Item price
				Total:        "10.00",
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Zero Price error\n",
		},
		{
			name: "Invalid Total Price",
			input: model.Receipt{
				Retailer:     "Walmart",
				PurchaseDate: "2023-08-10",
				PurchaseTime: "14:00",
				Items:        []model.Item{{"item1", "10.00"}},
				Total:        "0", // zero price
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Zero Price error\n",
		},
		{
			name: "Invalid Total Price",
			input: model.Receipt{
				Retailer:     "Walmart",
				PurchaseDate: "2023-08-10",
				PurchaseTime: "14:00",
				Items:        []model.Item{{"item1", "10.00"}},
				Total:        "-10.00", // negative price
			},
			httpStatus: http.StatusBadRequest,
			errorMsg:   "Invalid Price Format\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var reqBody []byte
			// decision to include the inputRaw struct field
			if tc.inputRaw != nil {
				reqBody = tc.inputRaw
			} else {
				reqBody, _ = json.Marshal(tc.input)
			}

			req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.ProcessReceipt(w, req)

			// Verify the HTTP status code
			if w.Code != tc.httpStatus {
				t.Errorf("Expected HTTP status code %d, got %d", tc.httpStatus, w.Code)
			}

			// Verify the error message
			if tc.errorMsg != "" {
				if w.Body.String() != tc.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, w.Body.String())
				}
			}
		})
	}
}

// test function for failure conditions in the GetPoints Endpoint
func TestGetPointsFailure_SingleCase(t *testing.T) {

	req := httptest.NewRequest("GET", "/getpoints/some-fake-id", nil)
	w := httptest.NewRecorder()

	// Expected HTTP status and error message
	expectedHttpStatus := http.StatusNotFound
	expectedErrorMsg := "No receipt found for that id\n"

	handler.GetPoints(w, req, "some-fake-id")

	// Verify the HTTP status code
	if w.Code != expectedHttpStatus {
		t.Errorf("Expected HTTP status code %d, got %d", expectedHttpStatus, w.Code)
	}

	// Verify the error message
	if w.Body.String() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, w.Body.String())
	}
}
