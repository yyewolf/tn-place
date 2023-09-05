package admin

import (
	"encoding/binary"
	"image"
	"image/draw"
	"os"

	"github.com/yyewolf/tn-place/back/internal/canva"
	"github.com/yyewolf/tn-place/back/internal/env"
	"github.com/yyewolf/tn-place/back/internal/server"

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

	img := server.Pl.Canva.Image
	// If the image is not the correct size, expand or crop it
	if img.Bounds().Dx() != r.Width || img.Bounds().Dy() != r.Height {
		newimg := image.NewNRGBA(image.Rect(0, 0, r.Width, r.Height))
		for i := range newimg.Pix {
			newimg.Pix[i] = 255
		}
		draw.Draw(newimg, newimg.Bounds(), img, image.Point{0, 0}, draw.Src)
		server.Pl.Canva.Image = newimg
		server.Pl.Imgbuf = nil
		os.WriteFile(env.C.SavePath, server.Pl.GetImageBytes(), 0644)

		newplacers := make([][]string, r.Width)
		for i := range newplacers {
			newplacers[i] = make([]string, r.Height)
		}
		for x := 0; x < img.Bounds().Dx(); x++ {
			for y := 0; y < img.Bounds().Dy(); y++ {
				if x >= r.Width || y >= r.Height {
					break
				}
				newplacers[x][y] = server.Pl.Canva.Placers[x][y]
			}
		}
		server.Pl.Canva.Placers = newplacers
		canva.SavePlacers(newplacers)
	}

	b := make([]byte, 32)
	binary.LittleEndian.PutUint32(b, uint32(52))
	server.Pl.Msgs <- b

	c.JSON(200, gin.H{"message": "ok"})
}
