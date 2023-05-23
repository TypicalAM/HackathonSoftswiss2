package routes

import (
	"net/http"
	"sort"

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

	var relations []models.ProfileTrash
	if err := con.db.Model(models.ProfileTrash{}).Where("profile_id = ?", user.Profile.ID).Find(&relations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot find your relations"})
		return
	}

	sort.Slice(relations, func(i, j int) bool {
		return relations[i].CreatedAt.After(relations[j].CreatedAt)
	})

	var trash []models.Product
	for _, relation := range relations {
		var product models.Product
		if err := con.db.Model(models.Product{}).Where("id = ?", relation.ProductID).First(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot find your trash"})
			return
		}

		trash = append(trash, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": user.Profile,
		"trash":   trash,
	})
}
