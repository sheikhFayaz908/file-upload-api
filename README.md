# File Upload API

This is a simple API for uploading and downloading CSV files in Go using the Gin framework and GORM for database operations.

## Features

- Upload a CSV file and store its data in the database.
- Download CSV data based on upload ID.
- Handle file processing asynchronously.
- Unit tests for upload and download functionality.

## Prerequisites

Before building or running this application, ensure you have the following prerequisites installed:
- C compiler (e.g., GCC) for CGO support 

## Requirements

- Go (version 1.13 or higher)
- Sqlite (or another relational database supported by GORM)
- Git


## RUNNING the Application

Follow these steps to RUN the application with CGO support:

1. Set the `CGO_ENABLED=1` environment variable to enable cgo support:
   export CGO_ENABLED=1
2. go run .
 
