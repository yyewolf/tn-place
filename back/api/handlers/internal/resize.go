package internal

import (
	"image"
	"image/draw"
	"tn-place/internal/server"

	"github.com/gin-gonic/gin"
)

type ResizeInput struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func resize(c *gin.Context) {
	var r ResizeInput
	err := c.BindJSON(&r)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Bad request"})
		return
	}

	img := server.Pl.Img
	// If the image is not the correct size, expand or crop it
	if img.Bounds().Dx() != r.Width || img.Bounds().Dy() != r.Height {
		newimg := image.NewNRGBA(image.Rect(0, 0, r.Width, r.Height))
		for i := range newimg.Pix {
			newimg.Pix[i] = 255
		}
		draw.Draw(newimg, newimg.Bounds(), img, image.Point{0, 0}, draw.Src)
		server.Pl.Img = newimg
	}
}
