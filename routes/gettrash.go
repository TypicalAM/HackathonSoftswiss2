package routes

import (
	"net/http"

	"github.com/TypicalAM/HackathonSoftswiss2/models"
	"github.com/gin-gonic/gin"
)

func (con controller) CheckEAN(c *gin.Context) {
	ean := c.Param("EAN")

	var product models.Product
	if res := con.db.Where("EAN = ?", ean).First(&product); res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}
