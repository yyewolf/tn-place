package image

import (
	"strconv"
	"tn-place/internal/server"

	"github.com/gin-gonic/gin"
)

func HandleGetImage(c *gin.Context) {
	b := server.Pl.GetImageBytes() //not thread safe but it won't do anything bad
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(b)))
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store")
	c.Writer.Write(b)
}

func LoadRoutes(r *gin.RouterGroup) {
	r.GET("/place.png", HandleGetImage)
}
