package databaser

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	dbhost, hok := os.LookupEnv("DB_HOST")
	if !hok {
		dbhost = "localhost"
	}
	dbuser, uok := os.LookupEnv("DB_USER")
	if !uok {
		dbuser = "postgres"
	}
	dbpass, pok := os.LookupEnv("DB_PASS")
	if !pok {
		log.Fatalln("no db pass")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s", dbhost, dbuser, dbpass)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("failed to start db %+v", err)
		os.Exit(1)
	}

	db.AutoMigrate(&AllowedUser{}, &IndexedArtist{}, &IndexedAlbum{}, &IndexedSong{})
	db.Raw("CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;")
	return db
}
