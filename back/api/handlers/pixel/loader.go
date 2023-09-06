package pixel

import (
	"fmt"
	"strconv"

	"github.com/yyewolf/tn-place/back/internal/server"

	"github.com/gin-gonic/gin"
)

func GetPixelInfo(c *gin.Context) {
	xP, yP := c.Param("x"), c.Param("y")
	// Parse x and y
	x, err := strconv.Atoi(xP)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Bad request"})
		return
	}
	y, err := strconv.Atoi(yP)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Bad request"})
		return
	}
	// Check if pixel is in bounds
	if x < 0 || x >= server.Pl.Canva.Image.Bounds().Dx() || y < 0 || y >= server.Pl.Canva.Image.Bounds().Dy() {
		c.AbortWithStatusJSON(400, gin.H{"error": "Bad request"})
		return
	}
	// Get pixel author
	author := server.Pl.Canva.Placers[x][y]

	if author == nil {
		c.JSON(200, gin.H{"placer": "Aucun"})
		return
	}

	c.JSON(200, gin.H{"placer": fmt.Sprintf("Équipe %s", author.Team)})
}

func LoadRoutes(r *gin.RouterGroup) {
	sg := r.Group("/pixel")
	sg.GET("/:x/:y/", GetPixelInfo)
}
