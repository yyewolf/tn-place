package server

import (
	"bytes"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"sync"
)

type Place struct {
	sync.RWMutex
	Msgs    chan []byte
	Close   chan int
	Clients []chan []byte
	Img     draw.Image
	Imgbuf  []byte
}

var Pl *Place

func NewServer(img draw.Image, count int) *Place {
	pl := &Place{
		RWMutex: sync.RWMutex{},
		Msgs:    make(chan []byte),
		Close:   make(chan int),
		Clients: make([]chan []byte, count),
		Img:     img,
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
		if err := png.Encode(buf, pl.Img); err != nil {
			log.Println(err)
		}
		pl.Imgbuf = buf.Bytes()
	}
	return pl.Imgbuf
}

func (pl *Place) SetPixel(x, y int, c color.Color) bool {
	rect := pl.Img.Bounds()
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	if 0 > x || x >= width || 0 > y || y >= height {
		return false
	}
	pl.Img.Set(x, y, c)
	pl.Imgbuf = nil
	return true
}