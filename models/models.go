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

// seedDatabase seeds the database with some data.
func seedDatabase(db *gorm.DB) error {
	products := []Product{
		{Name: "Butelka PET", EAN: "1000000002137"},
		{Name: "Dupa mojego starego", EAN: "2000000002137"},
		{Name: "Tytan gigant pierdolnik", EAN: "3000000002137"},
		{Name: "Jack Black", EAN: "4000000002137"},
	}

	var productsTwo []Product
	if res := db.Create(&products); res.Error != nil {
		log.Println("Failed to create products", res.Error)
		if res2 := db.Find(&productsTwo); res2.Error != nil {
			log.Println("Failed to find products", res2.Error)
			return res.Error
		}
	}

	// Create a user obj
	var user User
	if res := db.Preload("Profile").First(&user, "username = test"); res.Error != nil {
		log.Println("Failed to find user test", res.Error)
		return res.Error
	}

	if err := db.Model(&user.Profile).Association("Trash").Append(&productsTwo); err != nil {
		log.Println("Failed to append products to trash", err)
		return err
	}

	return nil
}

// Migrate migrates the database.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Product{}, &Session{}, &Profile{}); err != nil {
		return err
	}

	log.Println("seed status", seedDatabase(db))
	return nil
}
