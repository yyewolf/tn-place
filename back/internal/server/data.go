package server

import (
	"bytes"
	"image/color"
	"image/png"
	"log"
	"sync"
	"tn-place/internal/canva"
)

type Place struct {
	sync.RWMutex
	Msgs    chan []byte
	Close   chan int
	Clients []chan []byte
	Canva   *canva.Canva
	Imgbuf  []byte
}

var Pl *Place

func NewServer(cv *canva.Canva, count int) *Place {
	pl := &Place{
		RWMutex: sync.RWMutex{},
		Msgs:    make(chan []byte),
		Close:   make(chan int),
		Clients: make([]chan []byte, count),
		Canva:   cv,
	}
	go pl.broadcastLoop()
	return pl
}

func (pl *Place) ClientAmount() int {
	pl.RLock()
	defer pl.RUnlock()
	count := 0
	for _, ch := range pl.Clients {
		if ch != nil {
			count++
		}
	}
	return count
}

func (pl *Place) GetConnIndex() int {
	for i, client := range pl.Clients {
		if client == nil {
			return i
		}
	}
	return -1
}
func (pl *Place) broadcastLoop() {
	for {
		select {
		case i := <-pl.Close:
			if pl.Clients[i] != nil {
				close(pl.Clients[i])
				pl.Clients[i] = nil
			}
		case p := <-pl.Msgs:
			for i, ch := range pl.Clients {
				if ch != nil {
					select {
					case ch <- p:
					default:
						close(ch)
						pl.Clients[i] = nil
					}
				}
			}
		}
	}
}

func (pl *Place) GetImageBytes() []byte {
	if pl.Imgbuf == nil {
		buf := bytes.NewBuffer(nil)
		if err := png.Encode(buf, pl.Canva.Image); err != nil {
			log.Println(err)
		}
		pl.Imgbuf = buf.Bytes()
	}
	return pl.Imgbuf
}

var colorPalette = []color.RGBA{
	{0, 0, 0, 255},
	{0, 0, 255, 255},
	{0, 255, 0, 255},
	{0, 255, 255, 255},
	{255, 0, 0, 255},
	{255, 0, 255, 255},
	{255, 255, 0, 255},
	{255, 255, 255, 255},
	{128, 128, 128, 255},
	{0, 0, 128, 255},
	{0, 128, 0, 255},
	{0, 128, 128, 255},
	{128, 0, 0, 255},
	{128, 0, 128, 255},
	{128, 128, 0, 255},
	{192, 192, 192, 255},
}

func (pl *Place) SetPixel(x, y int, c color.Color) bool {
	// If the color is not in the 16-bit palette, return false.
	R, G, B, A := c.RGBA()
	for _, color := range colorPalette {
		if color.R == uint8(R>>8) && color.G == uint8(G>>8) && color.B == uint8(B>>8) && color.A == uint8(A>>8) {
			goto found
		}
	}
	return false
found:

	rect := pl.Canva.Image.Bounds()
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	if 0 > x || x >= width || 0 > y || y >= height {
		return false
	}
	pl.Canva.Image.Set(x, y, c)
	pl.Imgbuf = nil
	return true
}
