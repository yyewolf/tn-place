package admin

import (
	"github.com/yyewolf/tn-place/back/internal/server"

	"github.com/gin-gonic/gin"
)

func pause(c *gin.Context) {
	// GET should get the current pause state
	// POST should toggle the pause state
	if c.Request.Method == "GET" {
		c.JSON(200, gin.H{
			"paused": server.Pl.Paused,
		})
		return
	}
	server.Pl.Paused = !server.Pl.Paused
	c.JSON(200, gin.H{
		"paused": server.Pl.Paused,
	})
}
