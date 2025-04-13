package config

import (
	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Use SQLite instead of MySQL
)

var (
	db *gorm.DB
)

func Connect() {
	// d, err := gorm.Open("mysql", "root:shadman@tcp(127.0.0.1:3306)/bookstore?charset=utf8mb4&parseTime=True&loc=Local")

	// d, err := gorm.Open("mysql", "shadman:shadman@12@/simplerst?charset=utf8&parseTime=True&loc=Local")
	// d, err := gorm.Open("mysql", "shadman:Shadman@12@/simplerst?charset=utf8&parseTime=True&loc=Local")
	d, err := gorm.Open("sqlite3", "test.db") // SQLite will create a file `test.db`
	if err != nil {
		panic(err)
	}

	db = d
}

func GetDB() *gorm.DB {
	return db
}
