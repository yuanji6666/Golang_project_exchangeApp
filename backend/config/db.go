package config

import (
	"log"
	"time"
	"exchangeapp/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() {
	dsn := Appconfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to initialze database, got err : %v",err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("Failed to config database, got err : %v",err)
	}
	sqlDB.SetMaxIdleConns(Appconfig.Database.MaxIdlesConns)
	sqlDB.SetMaxOpenConns(Appconfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.Db =db
}
