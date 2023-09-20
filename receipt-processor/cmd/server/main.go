package main

import (
	"log"
	"net/http"
	"receipt-processor/internal/handler"
	"strings"
)

func main() {
	// Handles the "/receipts/process" route for processing receipts.
	// Accepts only POST requests with Content-Type "application/json".
	http.HandleFunc("/receipts/process", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content Type not allowed", http.StatusUnsupportedMediaType)
			return
		}

		handler.ProcessReceipt(w, r)
	})

	// Handles the "/receipts/" route for getting points associated with a receipt ID.
	// Accepts only GET requests.
	http.HandleFunc("/receipts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract the ID and optional "points" subpath from the URL
		pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/receipts/"), "/")
		id := pathSegments[0]
		if id == "" {
			http.Error(w, "Missing ID", http.StatusBadRequest)
			return
		}

		// If there's a "points" subpath, call GetPoints; otherwise do not allow
		if len(pathSegments) > 1 && pathSegments[1] == "points" {
			handler.GetPoints(w, r, id)
		} else {
			http.Error(w, "Method or Path not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
