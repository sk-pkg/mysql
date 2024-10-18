package main

import (
	"context"
	"github.com/sk-pkg/logger"
	"github.com/sk-pkg/mysql"
	"gorm.io/gorm"
	"log"
)

// Product represents a product in the database.
// It embeds gorm.Model which provides ID, CreatedAt, UpdatedAt, and DeletedAt fields.
type Product struct {
	gorm.Model
	Code  string // Product code
	Price uint   // Product price
}

// main is the entry point of the program.
// It demonstrates database connection, migration, and CRUD operations.
func main() {
	// Configure MySQL connection
	cfg := mysql.Config{
		User:     "homestead",
		Password: "secret",
		Host:     "127.0.0.1:33060",
		DBName:   "mysql_test",
	}

	// Initialize logger
	manager, err := logger.New()
	if err != nil {
		log.Fatal("failed to initialize logger:", err)
	}

	// Initialize database connection with custom logger
	db, err := mysql.New(mysql.WithConfigs(cfg), mysql.WithGormConfig(gorm.Config{
		Logger: mysql.NewLog(manager),
	}))
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Add trace_id to context for logging
	ctx := context.WithValue(context.Background(), "trace_id", "123456")
	db = db.WithContext(ctx)

	// Auto migrate the Product schema
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal("failed to auto migrate:", err)
	}

	// Perform CRUD operations

	// Create a new product
	createProduct(db)

	// Read products
	readProducts(db)

	// Update a product
	updateProduct(db)

	// Delete a product
	deleteProduct(db)
}

// createProduct demonstrates how to create a new product in the database.
//
// Parameters:
//   - db: *gorm.DB - The database connection
//
// Example:
//
//	createProduct(db)
func createProduct(db *gorm.DB) {
	product := &Product{Code: "D42", Price: 100}
	result := db.Create(product)
	if result.Error != nil {
		log.Printf("failed to create product: %v", result.Error)
	} else {
		log.Printf("created product with ID: %d", product.ID)
	}
}

// readProducts demonstrates how to read products from the database.
//
// Parameters:
//   - db: *gorm.DB - The database connection
//
// Example:
//
//	readProducts(db)
func readProducts(db *gorm.DB) {
	var product Product

	// Read the first product with ID 1
	result := db.First(&product, 1)
	if result.Error != nil {
		log.Printf("failed to read product with ID 1: %v", result.Error)
	} else {
		log.Printf("read product: %+v", product)
	}

	// Read the first product with code "D44"
	result = db.First(&product, "code = ?", "D44")
	if result.Error != nil {
		log.Printf("failed to read product with code D44: %v", result.Error)
	} else {
		log.Printf("read product: %+v", product)
	}
}

// updateProduct demonstrates how to update a product in the database.
//
// Parameters:
//   - db: *gorm.DB - The database connection
//
// Example:
//
//	updateProduct(db)
func updateProduct(db *gorm.DB) {
	var product Product

	// Update a single field
	result := db.Model(&product).Update("Price", 200)
	if result.Error != nil {
		log.Printf("failed to update product price: %v", result.Error)
	}

	// Update multiple fields using struct
	result = db.Model(&product).Updates(Product{Price: 200, Code: "F42"})
	if result.Error != nil {
		log.Printf("failed to update product using struct: %v", result.Error)
	}

	// Update multiple fields using map
	result = db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	if result.Error != nil {
		log.Printf("failed to update product using map: %v", result.Error)
	}
}

// deleteProduct demonstrates how to delete a product from the database.
//
// Parameters:
//   - db: *gorm.DB - The database connection
//
// Example:
//
//	deleteProduct(db)
func deleteProduct(db *gorm.DB) {
	result := db.Delete(&Product{}, 1) // Delete product with ID 1
	if result.Error != nil {
		log.Printf("failed to delete product: %v", result.Error)
	} else {
		log.Printf("deleted product with ID 1")
	}
}
