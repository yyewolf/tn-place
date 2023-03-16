package status

import (
	"github.com/yyewolf/tn-place/back/internal/server"

	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	count := server.Pl.ClientAmount()
	c.JSON(200, gin.H{
		"success": true,
		"status":  "OK",
		"clients": count,
	})
}
