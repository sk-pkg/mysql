# Mysql

Mysql是在gorm基础上封装的快速使用包。

## 安装

```shell
go get "github.com/sk-pkg/mysql"
```

## 快速开始
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

	// 迁移 schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // 根据整形主键查找
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(&product, 1)
}
```
## 可配置项
1.Configs数据库连接配置，如果要批量创建则填写多个配置，此项为必填
```go
mysql.WithConfigs(mysql.Config{
    User:     "homestead",
    Password: "secret",
    Host:     "127.0.0.1:33060",
    DBName:   "mysql_test",
})
```
多个数据库
```go
mysql.WithConfigs(config1, config2, config3)
```
2.gormConfig可自定义配置gorm，此项可选可参考[GORM 配置](https://gorm.io/zh_CN/docs/gorm_config.html)
```go
mysql.WithGormConfig(gorm.Config{})
```
3.maxIdleConn空闲连接池中连接的最大数量，默认值为10，此项可选
```go
mysql.WithMaxIdleConn(10)
```
4.maxOpenConn打开数据库连接的最大数量，默认值为50，此项可选
```go
mysql.WithMaxOpenConn(50)
```
5.connMaxLifetime连接可复用的最大时间，默认值为3小时，此项可选
```go
mysql.WithConnMaxLifetime(3 * time.Hour)
```
## 可用方法
1.New 初始化一个数据库连接实例
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
2.NewMulti 批量初始化数据库连接实例，成功会返回多个以数据库名为key，*gorm.DB为value的集合
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
## 参考文档
[GORM](https://gorm.io/zh_CN/)