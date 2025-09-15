package config

import (
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, models ...interface{}) {
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("Migration ล้มเหลว: %v", err)
	}
	log.Println("Auth Migration เสร็จสิ้น")
}
