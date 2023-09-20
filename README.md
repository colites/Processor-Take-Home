# receipt-processor

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
git clone https://github.com/colites/Processor-Take-Home.git
```

Navigate to the root directory:
```bash
cd Processor-Take-Home/receipt-processor
```

If not already initialized, run the following command to initialize the Go module
```bash
go mod init receipt-processor
```

Navigate to the project directory:
```bash
cd cmd/server
```
Build the application:
```bash
go build
```
Start the server
```bash
./server
```
or on windows
```bash
server.exe
```

The server will start running on http://localhost:8080.

## Usage

### Processing a receipt
To process a new receipt, make a POST request to /receipts/process.
on UNIX-like system:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"retailer":"Some Retailer","purchaseDate":"2023-09-18","purchaseTime":"15:04","items":[{"shortDescription":"item1","price":"10.00"},{"shortDescription":"item2","price":"20.00"}],"total":"30.00"}' http://localhost:8080/receipts/process
```

for windows:
```bash
curl -X POST -H "Content-Type: application/json" -d "{\"retailer\":\"Some Retailer\",\"purchaseDate\":\"2023-09-18\",\"purchaseTime\":\"15:04\",\"items\":[{\"shortDescription\":\"item1\",\"price\":\"10.00\"},{\"shortDescription\":\"item2\",\"price\":\"20.00\"}],\"total\":\"30.00\"}" http://localhost:8080/receipts/process
```

The API will return an ID for the stored receipt.

### Get Points for a receipt

To get the points for a processed receipt, make a GET request to /receipts/{id}.

For Unix systems:
```bash
curl --request GET \
  --url http://localhost:8080/receipts/{id}
```

for windows:
```bash
curl http://localhost:8080/receipts/{id}
```

The API will return the number of points for the given ID.

