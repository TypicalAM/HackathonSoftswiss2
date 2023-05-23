package models

import (
	"log"

	"gorm.io/gorm"
)

// User holds information about a user.
type User struct {
	gorm.Model
	Username string
	Password string
	Profile  Profile
	Sessions []Session
}

// Profile holds information about a user's profile.
type Profile struct {
	gorm.Model        `json:"-"`
	UserID            uint      `gorm:"unique" json:"-"`
	DisplayName       string    `json:"display_name"`
	ImageURL          string    `json:"image_url"`
	TotalSavedMass    int       `json:"total_saved_mass"`
	TotalPreventedCO2 int       `json:"total_prevented_co2"`
	Trash             []Product `gorm:"many2many:profile_trash" json:"-"`
}

// AfterCreate is a hook that is called to make sure that a profile is created for the user.
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	log.Println("Creating profile for user", u.Username)

	var count int64
	if res := tx.Model(&u.Profile).Where("user_id = ?", u.ID).Count(&count); res.Error != nil {
		return res.Error
	}

	if count > 0 {
		log.Println("User already has a profile")
		return
	}

	profile := Profile{
		UserID:      u.ID,
		DisplayName: u.Username,
		// TODO: Get default image URL from config
		ImageURL: "https://www.stockvault.net/data/2009/07/28/109653/preview16.jpg",
		TotalSavedMass: 0,
		TotalPreventedCO2: 0,
	}

	if res := tx.Create(&profile); res.Error != nil {
		log.Println("Error creating profile", res.Error)
		return res.Error
	}

	u.Profile = profile
	return nil
}
