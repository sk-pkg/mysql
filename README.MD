# sk-pkg/mysql Documentation

## Introduction

sk-pkg/mysql is a MySQL database operation package built on top of GORM, providing simple and easy-to-use database connection and operation interfaces, along with integrated custom logging functionality. This package aims to simplify MySQL database connection configuration and operations, allowing developers to focus more on implementing business logic.

## Installation

```bash
go get "github.com/sk-pkg/mysql"
```

## Basic Usage

### Single Database Connection

```go
import (
    "github.com/sk-pkg/mysql"
    "gorm.io/gorm"
    "log"
)

type Product struct {
    gorm.Model
    Code  string
    Price uint
}

func main() {
    // Configure database connection information
    cfg := mysql.Config{
        User:     "homestead",
        Password: "secret",
        Host:     "127.0.0.1:33060",
        DBName:   "mysql_test",
    }

    // Create database connection
    db, err := mysql.New(mysql.WithConfigs(cfg))
    if err != nil {
        log.Fatal("failed to connect database")
    }

    // Auto migrate table structure
    db.AutoMigrate(&Product{})

    // Create record
    db.Create(&Product{Code: "D42", Price: 100})

    // Query record
    var product Product
    db.First(&product, 1)                 // Query by primary key
    db.First(&product, "code = ?", "D42") // Query by condition

    // Update record
    db.Model(&product).Update("Price", 200)
    // Update multiple fields
    db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // Only update non-zero value fields
    db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

    // Delete record
    db.Delete(&product, 1)
}
```

### Multiple Database Connections

```go
import (
    "github.com/sk-pkg/mysql"
    "gorm.io/gorm"
    "log"
)

func main() {
    cfg1 := mysql.Config{
        User:     "homestead",
        Password: "secret",
        Host:     "127.0.0.1:33060",
        DBName:   "mysql_test",
    }

    cfg2 := mysql.Config{
        User:     "homestead",
        Password: "secret",
        Host:     "127.0.0.1:33060",
        DBName:   "mysql_test2",
    }

    // Create multiple database connections
    dbs, err := mysql.NewMulti(mysql.WithConfigs(cfg1, cfg2))
    if err != nil {
        log.Fatal("failed to connect database")
    }

    // Use specific database connection
    var product Product
    dbs["mysql_test"].First(&product, 1)
    
    var product2 Product
    dbs["mysql_test2"].First(&product2, 1)
}
```

## Configuration Options

sk-pkg/mysql provides various configuration options that can be set using the functional options pattern:

### Basic Configuration Options

1. **WithConfigs**: Set database connection configuration
   ```go
   cfg := mysql.Config{
       User:     "homestead",
       Password: "secret",
       Host:     "127.0.0.1:33060",
       DBName:   "mysql_test",
   }
   db, err := mysql.New(mysql.WithConfigs(cfg))
   ```

2. **WithGormConfig**: Set GORM configuration
   ```go
   gormCfg := gorm.Config{SkipDefaultTransaction: true}
   db, err := mysql.New(mysql.WithConfigs(cfg), mysql.WithGormConfig(gormCfg))
   ```

3. **WithMaxIdleConn**: Set maximum idle connections
   ```go
   db, err := mysql.New(mysql.WithConfigs(cfg), mysql.WithMaxIdleConn(20))
   ```

4. **WithMaxOpenConn**: Set maximum open connections
   ```go
   db, err := mysql.New(mysql.WithConfigs(cfg), mysql.WithMaxOpenConn(100))
   ```

5. **WithConnMaxLifetime**: Set maximum connection lifetime
   ```go
   db, err := mysql.New(mysql.WithConfigs(cfg), mysql.WithConnMaxLifetime(4 * time.Hour))
   ```

## Logging Functionality

sk-pkg/mysql integrates custom logging functionality to record SQL queries and execution details.

### Installing Dependencies

```bash
go get github.com/sk-pkg/logger
go get gorm.io/gorm
```

### Creating a Logger

```go
import (
    "github.com/sk-pkg/logger"
    "github.com/sk-pkg/mysql"
    "gorm.io/gorm"
)

// Initialize logger manager
manager, err := logger.New()
if err != nil {
    log.Fatal("failed to initialize logger:", err)
}

// Create logger
logger := mysql.NewLog(manager)

// Use logger in GORM configuration
db, err := mysql.New(
    mysql.WithConfigs(cfg), 
    mysql.WithGormConfig(gorm.Config{
        Logger: logger,
    })
)
```

### Logger Configuration Options

1. **WithLevel**: Set log level
   ```go
   logger := mysql.NewLog(manager, mysql.WithLevel("info"))
   ```
   Supported log levels:
   - "info" (default): Records all SQL queries
   - "warn": Records only warnings and errors
   - "error": Records only errors
   - "silent": Records no logs

2. **WithSlowThreshold**: Set slow query threshold
   ```go
   logger := mysql.NewLog(manager, mysql.WithSlowThreshold(300 * time.Millisecond))
   ```

3. **WithIgnoreRecordNotFoundError**: Whether to ignore record not found errors
   ```go
   logger := mysql.NewLog(manager, mysql.WithIgnoreRecordNotFoundError(true))
   ```

4. **Combining multiple options**:
   ```go
   logger := mysql.NewLog(manager,
       mysql.WithLevel("info"),
       mysql.WithIgnoreRecordNotFoundError(true),
       mysql.WithSlowThreshold(300 * time.Millisecond),
   )
   ```

### Using Context for Log Tracing

```go
// Add trace_id to context for log tracing
ctx := context.WithValue(context.Background(), "trace_id", "123456")
db = db.WithContext(ctx)

// Execute database operations, logs will include trace_id
db.Create(&Product{Code: "D42", Price: 100})
```

## Complete Example

```go
package main

import (
    "context"
    "github.com/sk-pkg/logger"
    "github.com/sk-pkg/mysql"
    "gorm.io/gorm"
    "log"
    "time"
)

type Product struct {
    gorm.Model
    Code  string
    Price uint
}

func main() {
    // Configure database connection
    cfg := mysql.Config{
        User:     "homestead",
        Password: "secret",
        Host:     "127.0.0.1:33060",
        DBName:   "mysql_test",
    }

    // Initialize logger manager
    manager, err := logger.New()
    if err != nil {
        log.Fatal("failed to initialize logger:", err)
    }

    // Create logger
    logger := mysql.NewLog(manager,
        mysql.WithLevel("info"),
        mysql.WithIgnoreRecordNotFoundError(true),
        mysql.WithSlowThreshold(300 * time.Millisecond),
    )

    // Initialize database connection
    db, err := mysql.New(
        mysql.WithConfigs(cfg),
        mysql.WithGormConfig(gorm.Config{
            Logger: logger,
        }),
        mysql.WithMaxIdleConn(10),
        mysql.WithMaxOpenConn(50),
        mysql.WithConnMaxLifetime(3 * time.Hour),
    )
    if err != nil {
        log.Fatal("failed to connect database:", err)
    }

    // Add trace_id to context for log tracing
    ctx := context.WithValue(context.Background(), "trace_id", "123456")
    db = db.WithContext(ctx)

    // Auto migrate table structure
    err = db.AutoMigrate(&Product{})
    if err != nil {
        log.Fatal("failed to auto migrate:", err)
    }

    // Create record
    product := &Product{Code: "D42", Price: 100}
    result := db.Create(product)
    if result.Error != nil {
        log.Printf("failed to create product: %v", result.Error)
    }

    // Query record
    var foundProduct Product
    result = db.First(&foundProduct, product.ID)
    if result.Error != nil {
        log.Printf("failed to read product: %v", result.Error)
    } else {
        log.Printf("read product: %+v", foundProduct)
    }

    // Update record
    result = db.Model(&foundProduct).Update("Price", 200)
    if result.Error != nil {
        log.Printf("failed to update product: %v", result.Error)
    }

    // Delete record
    result = db.Delete(&foundProduct)
    if result.Error != nil {
        log.Printf("failed to delete product: %v", result.Error)
    }
}
```

## Considerations

1. In production environments, set appropriate log levels to avoid excessive log output that could impact performance.
2. The slow query threshold should be set based on your application requirements and database performance.
3. Using `WithIgnoreRecordNotFoundError(true)` can reduce unnecessary "record not found" error logs, but ensure this doesn't affect your error handling logic.
4. Default connection pool configuration:
   - Maximum idle connections: 10
   - Maximum open connections: 50
   - Maximum connection lifetime: 3 hours

## Reference Documentation

For more information about using GORM, please refer to the [GORM Official Documentation](https://gorm.io/en/).