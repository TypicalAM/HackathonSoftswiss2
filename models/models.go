package models

import (
	"log"

	"github.com/TypicalAM/HackathonSoftswiss2/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New connects to the database using the config.
func New(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate migrates the database.
func Migrate(db *gorm.DB) error {
	if err := db.SetupJoinTable(&Profile{}, "Trash", &ProfileTrash{}); err != nil {
		return err
	}


	if err := db.AutoMigrate(&User{}, &Product{}, &Session{}, &Profile{}); err != nil {
		return err
	}

	log.Println("seed status", seedDatabase(db))
	return nil
}
