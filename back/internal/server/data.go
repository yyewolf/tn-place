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

// colorPalette = ['#F2F3F4', '#222222', '#F3C300', '#875692', '#F38400', '#A1CAF1', '#BE0032', '#C2B280', '#848482', '#008856', '#E68FAC', '#0067A5', '#F99379', '#604E97', '#F6A600', '#B3446C', '#DCD300', '#882D17', '#8DB600', '#654522', '#E25822', '#2B3D26'];

// var colorPalette = []color.Color{
// 	color.RGBA{242, 243, 244, 255},
// 	color.RGBA{34, 34, 34, 255},
// 	color.RGBA{243, 195, 0, 255},
// 	color.RGBA{135, 86, 146, 255},
// 	color.RGBA{243, 132, 0, 255},
// 	color.RGBA{161, 202, 241, 255},
// 	color.RGBA{190, 0, 50, 255},
// 	color.RGBA{194, 178, 128, 255},
// 	color.RGBA{132, 132, 130, 255},
// 	color.RGBA{0, 136, 86, 255},
// 	color.RGBA{230, 143, 172, 255},
// 	color.RGBA{0, 103, 165, 255},
// 	color.RGBA{249, 147, 121, 255},
// 	color.RGBA{96, 78, 151, 255},
// 	color.RGBA{246, 166, 0, 255},
// 	color.RGBA{179, 68, 108, 255},
// 	color.RGBA{220, 211, 0, 255},
// 	color.RGBA{136, 45, 23, 255},
// 	color.RGBA{141, 182, 0, 255},
// 	color.RGBA{101, 69, 34, 255},
// 	color.RGBA{226, 88, 34, 255},
// 	color.RGBA{43, 61, 38, 255},
// }

func (pl *Place) SetPixel(x, y int, c color.Color) bool {
	// If the color is not in the 16-bit palette, return false.
	// 	r, g, b, a := c.RGBA()
	// 	for _, color := range colorPalette {
	// 		R, G, B, A := color.RGBA()
	// 		if r == R && g == G && b == B && a == A {
	// 			goto found
	// 		}
	// 	}
	// 	return false
	// found:

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
