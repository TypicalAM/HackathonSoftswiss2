package routes

import (
	"net/http"

	"github.com/TypicalAM/HackathonSoftswiss2/models"
	"github.com/gin-gonic/gin"
)

// Profile fetches the user's profile.
func (con controller) Profile(c *gin.Context) {
	user, err := con.getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	var trash []models.Product
	if err := con.db.Model(&user.Profile).Association("Trash").Find(&trash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot find your trash"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": user.Profile,
		"trash":   trash,
	})
}
