package config

import (
	"fmt"
	"log"
	"orders/model"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db  *gorm.DB
	err error
)

func ConnectDB() {
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&interpolateParams=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("ORDER_DB_DATABASE"))
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Database connection failed: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(150)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to database")
}

func GetDB() *gorm.DB {
	if os.Getenv("ENV") == "local" {
		return db.Debug()
	}

	return db
}

func SyncDB() {
	db.AutoMigrate(&model.Cart{})
	db.AutoMigrate(&model.CartItem{})
	db.AutoMigrate(&model.Order{})
	db.AutoMigrate(&model.OrderItem{})
}
