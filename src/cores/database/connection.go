package database

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() *gorm.DB {
	// Loan ENV
	err := godotenv.Load()

	// Check if there is error when load ENV
	if err != nil {
		log.Fatal("Failed to Load .env")
	}

	// Assign ENV value to Variable
	dbUsername := os.Getenv("MYSQL_USERNAME")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DBNAME")

	var dsn string
	if dbPassword == "" {
		dsn = fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbHost, dbPort, dbName)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	// Check if there is error when connect to database
	if err != nil {
		log.Fatal("Failed connect to database: ", err.Error())
	}

	// Double Check Database Connection ( As Recommended by Official Documentation )
	// if errorPing := db.Ping(); err != nil {
	// 	log.Fatal("Failed connect to database: ", errorPing.Error())
	// }

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)

	DB = db
	return db
}
