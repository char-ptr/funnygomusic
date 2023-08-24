package databaser

import (
	"fmt"
	"gorm.io/driver/postgres"
	"log"
	"os"

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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=music", dbhost, dbuser, dbpass)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln("failed to start db", err)
	}

	db.AutoMigrate(&AllowedUser{}, &IndexedArtist{}, &IndexedAlbum{}, &IndexedSong{})
	db.Raw("CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;")
	return db
}
