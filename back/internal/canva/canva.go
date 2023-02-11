package canva

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"tn-place/internal/env"
)

func NewImage() draw.Image {
	var img draw.Image
	if env.LoadPath == "" {
		nrgba := image.NewNRGBA(image.Rect(0, 0, env.Width, env.Height))
		for i := range nrgba.Pix {
			nrgba.Pix[i] = 255
		}
		img = nrgba
	} else {
		img = Load(env.LoadPath)

		// If the image is not the correct size, expand it only
		if img.Bounds().Dx() < env.Width || img.Bounds().Dy() < env.Height {
			newimg := image.NewNRGBA(image.Rect(0, 0, env.Width, env.Height))
			for i := range newimg.Pix {
				newimg.Pix[i] = 255
			}
			draw.Draw(newimg, newimg.Bounds(), img, image.Point{0, 0}, draw.Src)
			img = newimg
		}
	}
	return img
}

func Load(loadPath string) draw.Image {
	f, err := os.Open(loadPath)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	pngimg, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	return pngimg.(draw.Image)
}
