# Mysql

Mysql is a quick-use package encapsulated on top of GORM.

## Installation

```shell
go get "github.com/sk-pkg/mysql"
```

## Quick Start

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
	cfg := mysql.Config{
		User:     "homestead",
		Password: "secret",
		Host:     "127.0.0.1:33060",
		DBName:   "mysql_test",
	}

	db, err := mysql.New(mysql.WithConfigs(cfg))
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // Find by integer primary key
	db.First(&product, "code = ?", "D42") // Find record where code field value is D42

	// Update - Set product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - Update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // Only update non-zero value fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - Delete product
	db.Delete(&product, 1)
}
```

## Configurable Options

1. **Configs**: Database connection configuration. If you want to create multiple instances, fill in multiple configurations. This is required.
```go
mysql.WithConfigs(mysql.Config{
    User:     "homestead",
    Password: "secret",
    Host:     "127.0.0.1:33060",
    DBName:   "mysql_test",
})
```

Multiple databases:
```go
mysql.WithConfigs(config1, config2, config3)
```

2. **gormConfig**: You can customize GORM configuration. This is optional. Refer to [GORM Configuration](https://gorm.io/en/docs/gorm_config.html).
```go
mysql.WithGormConfig(gorm.Config{})
```

3. **maxIdleConn**: The maximum number of connections in the idle connection pool, with a default value of 10. This is optional.
```go
mysql.WithMaxIdleConn(10)
```

4. **maxOpenConn**: The maximum number of open database connections, with a default value of 50. This is optional.
```go
mysql.WithMaxOpenConn(50)
```

5. **connMaxLifetime**: The maximum lifetime of a reusable connection, with a default value of 3 hours. This is optional.
```go
mysql.WithConnMaxLifetime(3 * time.Hour)
```

## Available Methods

1. **New**: Initialize a database connection instance.
```go
cfg := mysql.Config{
		User:     "homestead",
		Password: "secret",
		Host:     "127.0.0.1:33060",
		DBName:   "mysql_test",
	}

	db, err := mysql.New(mysql.WithConfigs(cfg))
	if err != nil {
		log.Fatal("failed to connect database")
	}
```

2. **NewMulti**: Batch initialize database connection instances. It will return a collection with database names as keys and *gorm.DB as values.
```go
cfg1 := Config{
	User:     "homestead",
	Password: "secret",
	Host:     "127.0.0.1:33060",
	DBName:   "mysql_test1",
}

cfg2 := Config{
	User:     "homestead",
	Password: "secret",
	Host:     "127.0.0.1:33060",
	DBName:   "mysql_test2",
}
    
dbs, err := NewMulti(WithConfigs(cfg1, cfg2))
if err != nil {
	t.Fatal("failed to connect database", err)
}
        
mysql1 := dbs["mysql_test1"]
mysql2 := dbs["mysql_test2"]
```

## Reference Document
[GORM](https://gorm.io/en/)