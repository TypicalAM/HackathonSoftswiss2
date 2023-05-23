package routes

import (
	"net/http"

	"github.com/TypicalAM/HackathonSoftswiss2/models"
	"github.com/gin-gonic/gin"
)

func (con controller) GlobalInfo(c *gin.Context) {
	// VERY Computationally expensive operation, but cmon, it's a hackathon
	var users []models.User
	if res := con.db.Preload("Profile").Find(&users); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get users"})
		return
	}

	var totalSavedMass int
	var totalPreventedCO2 int

	var products []models.Product
	for _, user := range users {
		if err := con.db.Model(&user.Profile).Association("Trash").Find(&products); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get products"})
			return
		}

		for _, product := range products {
			totalSavedMass += product.Mass
			totalPreventedCO2 += product.CO2EmissionPrevented
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_saved_mass":    totalSavedMass,
		"total_prevented_co2": totalPreventedCO2,
	})
}
