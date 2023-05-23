package routes

import (
	"log"
	"net/http"

	"github.com/TypicalAM/HackathonSoftswiss2/models"
	"github.com/gin-gonic/gin"
)

type ThrowAwayData struct {
	EAN string `json:"ean"`
}

func (con controller) ThrowAway(c *gin.Context) {
	var tad ThrowAwayData
	if err := c.ShouldBindJSON(&tad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if res := con.db.Where("ean = ?", tad.EAN).First(&product); res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	user, err := con.getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	if err := con.db.Model(&models.ProfileTrash{}).Create(&models.ProfileTrash{
		ProfileID: user.Profile.ID,
		ProductID: product.ID,
	}); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create relation"})
		log.Println(err)
		return
	}

	if res := con.db.Model(&user.Profile).Update("total_saved_mass", user.Profile.TotalSavedMass+product.Mass); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update total saved mass"})
		log.Println(err)
		return
	}

	if res := con.db.Model(&user.Profile).Update("total_prevented_co2", user.Profile.TotalPreventedCO2+product.CO2EmissionPrevented); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update total saved co2"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trash thrown"})
}
