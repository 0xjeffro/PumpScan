package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var conn *gorm.DB

func InitDB() error {

	dsn := func() string {
		if os.Getenv("DSN") == "" {
			panic("DSN is not set!")
		}
		return os.Getenv("DSN")
	}()

	var err error
	conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("[GetConn] Gorm open error: ", err)
		return err
	}
	log.Println("[GetConn] Gorm open success")

	sqlDB, err := conn.DB()
	if err != nil {
		log.Println("[GetConn] Get DB error: ", err)
		return err
	}
	sqlDB.SetMaxIdleConns(30)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return nil
}

func GetConn() (*gorm.DB, error) {
	if conn == nil {
		err := InitDB()
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}
