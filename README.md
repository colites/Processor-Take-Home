# Processor-Take-Home

The Receipt Points Processor is a Go-based web application that provides APIs for storing receipts and calculating points based on given rules. This application uses in-memory storage to manage receipt information.

## Features
Store a receipt and receive an ID for future reference

Calculate points for each stored receipt based on a set of rules

Retrieve points for a stored receipt by ID

## Prerequisites
Go 1.16+ (for running the server)

A tool to make HTTP requests (e.g., curl, Postman)

## Installation

Clone the repository:
```bash
git clone https://github.com/colites/receipt-processor.git
```
Navigate to the project directory:
```bash
cd receipt-processor
```
Build and run the application:
```bash
go build
./receipt-processor

```

The server will start running on http://localhost:8080.

## Usage

### Processing a receipt
To process a new receipt, make a POST request to /receipts/process.
```bash
curl --request POST \
  --url http://localhost:8080/receipts/process \
  --header 'Content-Type: application/json' \
  --data '{
    "retailer": "SomeNewWalmart",
    "purchaseDate": "2023-09-20",
    "purchaseTime": "14:00",
    "total": "10.25",
    "items": [
      {
        "shortDescription": "Milk",
        "price": "2.50"
      },
      {
        "shortDescription": "Candy",
        "price": "1.75"
      }
    ]
  }'
```
The API will return an ID for the stored receipt.

### Get Points for a receipt

To get the points for a processed receipt, make a GET request to /receipts/{id}.
```bash
curl --request GET \
  --url http://localhost:8080/receipts/{id}
```

The API will return the number of points for the given ID.

