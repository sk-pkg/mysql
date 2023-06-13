package mysql

import (
	"fmt"
	"gorm.io/gorm"
	"testing"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestMysql(t *testing.T) {
	cfg := Config{
		User:     "homestead",
		Password: "secret",
		Host:     "127.0.0.1:33060",
		DBName:   "mysql_test",
	}

	db, err := New(WithConfigs(cfg))
	if err != nil {
		t.Fatal("failed to connect database", err)
	}

	var product Product
	err = db.First(&product, 1).Error
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Test Mysql Success output:%v", product)
}

func TestMysqlMulti(t *testing.T) {
	cfg1 := Config{
		User:     "homestead",
		Password: "secret",
		Host:     "127.0.0.1:33060",
		DBName:   "mysql_test",
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

	var product, product2 Product
	err = dbs["mysql_test"].First(&product, 1).Error
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Test mysql_test Success output:%v \n", product)

	err = dbs["mysql_test2"].First(&product2, 1).Error
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Test mysql_test2 Success output:%v", product2)
}
