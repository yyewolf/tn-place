package canva

import (
	"encoding/json"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"tn-place/internal/env"
)

var PlacerFile string

type Canva struct {
	Width   int
	Height  int
	Image   draw.Image
	Placers [][]string
}

func NewImage() *Canva {
	var img draw.Image
	var placers [][]string
	if env.LoadPath == "" {
		nrgba := image.NewNRGBA(image.Rect(0, 0, env.Width, env.Height))
		for i := range nrgba.Pix {
			nrgba.Pix[i] = 255
		}
		img = nrgba

		placers = make([][]string, env.Width)
		for i := range placers {
			placers[i] = make([]string, env.Height)
		}
	} else {
		img = Load(env.LoadPath)
		if img == nil {
			nrgba := image.NewNRGBA(image.Rect(0, 0, env.Width, env.Height))
			for i := range nrgba.Pix {
				nrgba.Pix[i] = 255
			}
			img = nrgba

			placers = make([][]string, env.Width)
			for i := range placers {
				placers[i] = make([]string, env.Height)
			}

			goto end
		}

		// Import placers from the same file but with a .json extension

		file := env.LoadPath
		if file[len(file)-4:] == ".png" {
			file = file[:len(file)-4]
		}
		file += ".json"
		PlacerFile = file
		placers = LoadPlacers(file)
		if placers == nil {
			placers = make([][]string, env.Width)
			for i := range placers {
				placers[i] = make([]string, env.Height)
			}
		}

		// If the image is not the correct size, expand it only
		if img.Bounds().Dx() < env.Width || img.Bounds().Dy() < env.Height {
			newimg := image.NewNRGBA(image.Rect(0, 0, env.Width, env.Height))
			for i := range newimg.Pix {
				newimg.Pix[i] = 255
			}
			draw.Draw(newimg, newimg.Bounds(), img, image.Point{0, 0}, draw.Src)
			img = newimg

			// Expand the placers
			newplacers := make([][]string, env.Width)
			for i := range newplacers {
				newplacers[i] = make([]string, env.Height)
			}
			for x := 0; x < img.Bounds().Dx(); x++ {
				for y := 0; y < img.Bounds().Dy(); y++ {
					newplacers[x][y] = placers[x][y]
				}
			}
		}
	}
end:
	return &Canva{
		Width:   env.Width,
		Height:  env.Height,
		Image:   img,
		Placers: placers,
	}
}

func Load(loadPath string) draw.Image {
	f, err := os.Open(loadPath)
	defer f.Close()
	if err != nil {
		return nil
	}
	pngimg, err := png.Decode(f)
	if err != nil {
		return nil
	}
	return pngimg.(draw.Image)
}

func LoadPlacers(loadPath string) [][]string {
	f, err := os.Open(loadPath)
	defer f.Close()
	if err != nil {
		return nil
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}
	var placers [][]string
	err = json.Unmarshal(data, &placers)
	if err != nil {
		return nil
	}
	return placers
}

func SavePlacers(placers [][]string) {
	f, err := os.Create(PlacerFile)
	defer f.Close()
	if err != nil {
		return
	}
	// Write the placers to the file
	data, err := json.Marshal(placers)
	if err != nil {
		return
	}
	f.Write(data)
}
