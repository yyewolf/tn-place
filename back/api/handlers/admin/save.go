package admin

import (
	"os"

	"github.com/yyewolf/tn-place/back/internal/canva"
	"github.com/yyewolf/tn-place/back/internal/env"
	"github.com/yyewolf/tn-place/back/internal/server"

	"github.com/gin-gonic/gin"
)

func save(c *gin.Context) {
	os.WriteFile(env.C.SavePath, server.Pl.GetImageBytes(), 0644)
	canva.SavePlacers(server.Pl.Canva.Placers)
	c.JSON(200, gin.H{
		"message": "Saved",
	})
}
