package routes

import (
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

	if err := con.db.Model(&user.Profile).Association("Trash").Append(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot throw trash"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trash thrown"})
}