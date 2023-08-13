package databaser

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:botte.db?cache=shared&mode=rwc"), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to start db", err)
	}

	db.AutoMigrate(&IndexEntry{}, &AllowedUser{})

	return db
}
