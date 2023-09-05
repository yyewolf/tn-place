package canva

import (
	"encoding/json"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"

	"github.com/yyewolf/tn-place/back/internal/env"
)

var PlacerFile string

type PlacerInfo struct {
	Name string
	Team int
}

type Canva struct {
	Width   int
	Height  int
	Image   draw.Image
	Placers [][]*PlacerInfo
}

func CreateRawImage(width, height int) (img draw.Image, placers [][]*PlacerInfo) {
	nrgba := image.NewNRGBA(image.Rect(0, 0, env.C.Width, env.C.Height))
	for i := range nrgba.Pix {
		nrgba.Pix[i] = 255
	}
	img = nrgba

	placers = make([][]*PlacerInfo, env.C.Width)
	for i := range placers {
		placers[i] = make([]*PlacerInfo, env.C.Height)
	}

	return img, placers
}

func NewImage() *Canva {
	var img draw.Image
	var placers [][]*PlacerInfo
	canva := &Canva{
		Width:   env.C.Width,
		Height:  env.C.Height,
		Image:   img,
		Placers: placers,
	}
	if env.C.LoadPath == "" {
		img, placers = CreateRawImage(env.C.Width, env.C.Height)
		canva.Image = img
		canva.Placers = placers
		return canva
	}

	img = Load(env.C.LoadPath)
	if img == nil {
		img, placers = CreateRawImage(env.C.Width, env.C.Height)
		canva.Image = img
		canva.Placers = placers
		return canva
	}

	// Import placers from the same file but with a .json extension
	file := env.C.LoadPath
	if file[len(file)-4:] == ".png" {
		file = file[:len(file)-4]
	}
	file += ".json"
	PlacerFile = file
	placers = LoadPlacers(file)
	if placers == nil {
		placers = make([][]*PlacerInfo, env.C.Width)
		for i := range placers {
			placers[i] = make([]*PlacerInfo, env.C.Height)
		}
	}

	// If the image is not the correct size, expand it only
	if img.Bounds().Dx() < env.C.Width || img.Bounds().Dy() < env.C.Height {
		img, placers = ExpandImage(img, placers, env.C.Width, env.C.Height)
	}

	canva.Image = img
	canva.Placers = placers
	return canva
}

func ExpandImage(img draw.Image, placers [][]*PlacerInfo, width, height int) (draw.Image, [][]*PlacerInfo) {
	newimg := image.NewNRGBA(image.Rect(0, 0, env.C.Width, env.C.Height))
	for i := range newimg.Pix {
		newimg.Pix[i] = 255
	}
	draw.Draw(newimg, newimg.Bounds(), img, image.Point{0, 0}, draw.Src)
	img = newimg

	// Expand the placers
	newplacers := make([][]*PlacerInfo, env.C.Width)
	for i := range newplacers {
		newplacers[i] = make([]*PlacerInfo, env.C.Height)
	}
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			newplacers[x][y] = placers[x][y]
		}
	}
	return img, placers
}

func Load(loadPath string) draw.Image {
	f, err := os.Open(loadPath)
	if err != nil {
		return nil
	}
	defer f.Close()
	pngimg, err := png.Decode(f)
	if err != nil {
		return nil
	}
	return pngimg.(draw.Image)
}

func LoadPlacers(loadPath string) [][]*PlacerInfo {
	f, err := os.Open(loadPath)
	if err != nil {
		return nil
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil
	}
	var placers [][]*PlacerInfo
	err = json.Unmarshal(data, &placers)
	if err != nil {
		return nil
	}
	return placers
}

func SavePlacers(placers [][]*PlacerInfo) {
	f, err := os.Create(PlacerFile)
	if err != nil {
		return
	}
	defer f.Close()
	// Write the placers to the file
	data, err := json.Marshal(placers)
	if err != nil {
		return
	}
	f.Write(data)
}
