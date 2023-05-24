package models

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type seedData []struct {
	Name     string       `json:"name"`
	ImageURL string       `json:"image_url"`
	CO2      int          `json:"co2"`
	Mass     int          `json:"mass"`
	EAN      string       `json:"ean"`
	BinType  TrashBinType `json:"type_of_trash"`
}

// seedDatabase seeds the database with some data.
func seedDatabase(db *gorm.DB) error {
	// Create some users
	if err := createUsers(db); err != nil {
		return err
	}

	// Create some products
	if err := createProducts(db); err != nil {
		return err
	}

	// Create some associations
	if err := createAssociations(db); err != nil {
		return err
	}

	return nil
}

// createUsers creates some users.
func createUsers(db *gorm.DB) error {
	log.Println("Creating users...")
	file, err := os.Open(filepath.Join("data", "users.txt"))
	if err != nil {
		return err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	var users []User
	for scan.Scan() {
		users = append(users, User{Username: scan.Text(), Password: "nie_zahashowany_okoń"})
	}

	var count int64
	if err := db.Model(User{}).Where("username = ?", users[0].Username).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		if err := db.Find(&users).Error; err != nil {
			return err
		}

		return nil
	}

	if err := db.Create(&users).Error; err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return err
	}

	return nil
}

// createProducts creates all products.
func createProducts(db *gorm.DB) error {
	log.Println("Creating products...")
	file, err := os.Open(filepath.Join("data", "products.json"))
	if err != nil {
		return err
	}
	defer file.Close()

	var data seedData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	var count int64
	if err := db.Model(Product{}).Where("name = ?", data[0].Name).Count(&count).Error; err != nil {
		return err
	}

	var products []Product
	if count > 0 {
		if err := db.Find(&products).Error; err != nil {
			return err
		}

		return nil
	}

	for _, product := range data {
		products = append(products, Product{
			Name:                 product.Name,
			CO2EmissionPrevented: product.CO2,
			Mass:                 product.Mass,
			EAN:                  product.EAN,
			TypeOfTrash:          product.BinType,
			ImageURL:             product.ImageURL,
		})
	}

	if err := db.Create(&products).Error; err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return err
	}

	return nil
}

// createAssociations creates all associations.
func createAssociations(db *gorm.DB) error {
	log.Println("Creating associations...")
	rand.Seed(time.Now().UnixNano())
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	var user User
	if err := db.Model(&user).Preload("Profile").Where("username = ?", "test").First(&user).Error; err != nil {
		return err
	}

	if db.Model(&user.Profile).Association("Trash").Count() > 0 {
		log.Println("num of trash for user", user.Username, ":", db.Model(&user.Profile).Association("Trash").Count())
		return nil
	}

	log.Println("num of trash for user", user.Username, ":", 1)
	numOfTrash := rand.Intn(10) + 1
	for i := 0; i < numOfTrash; i++ {
		// It is slow, but it works.
		num := rand.Intn(len(products))
		if err := db.Model(&user.Profile).Association("Trash").Append(&products[num]); err != nil {
			return err
		}

		if res := db.Model(&user.Profile).Update("total_saved_mass", user.Profile.TotalSavedMass+products[num].Mass); res.Error != nil {
			return res.Error
		}

		if res := db.Model(&user.Profile).Update("total_prevented_co2", user.Profile.TotalPreventedCO2+products[num].CO2EmissionPrevented); res.Error != nil {
			return res.Error
		}

	}

	var users []User
	if err := db.Preload("Profile").Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		if db.Model(&user.Profile).Association("Trash").Count() > 0 {
			return nil
		}

		numOfTrash := rand.Intn(10) + 1
		log.Println("num of trash for user", user.Username, ":", numOfTrash)
		for i := 0; i < numOfTrash; i++ {
			// It is slow, but it works.
			if err := db.Model(&user.Profile).Association("Trash").Append(&products[rand.Intn(len(products))]); err != nil {
				return err
			}
		}
	}

	return nil
}
