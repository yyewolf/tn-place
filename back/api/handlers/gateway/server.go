package gateway

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"time"
	"tn-place/internal/env"
	"tn-place/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  32,
	WriteBufferSize: 32,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, req *http.Request, status int, err error) {
		log.Println(err)
		http.Error(w, "Error while trying to make websocket connection.", status)
	},
}

func GetGateway(c *gin.Context) {
	server.Pl.Lock()
	defer server.Pl.Unlock()
	i := server.Pl.GetConnIndex()
	if i == -1 {
		log.Println("Server full.")
		http.Error(c.Writer, "Server full.", 503)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	ch := make(chan []byte, 8)
	server.Pl.Clients[i] = ch

	go readLoop(conn, i, c, ch)
	go writeLoop(conn, ch)

	waiterID := c.MustGet("user_id").(string)

	// Send timer to client
	w, ok := waiter[waiterID]
	if !ok {
		return
	}
	s := w.Sub(time.Now()).Seconds()
	if s < 0 {
		s = 0
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(s))
	ch <- b
}

var waiter = make(map[string]time.Time)

func readLoop(conn *websocket.Conn, i int, c *gin.Context, ch chan []byte) {
	waiterID := c.MustGet("user_id").(string)
	limiter := rateLimiter()
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		if !limiter() {
			log.Printf("[ERR] %s got recked by rate limiter.\n", waiterID)
			break
		}
		w, ok := waiter[waiterID]
		if ok && w.After(time.Now()) {
			log.Printf("[ERR] %s was ignored due to timeout.\n", waiterID)
			continue
		}
		if messageHandler(p) != nil {
			log.Printf("[ERR] %s sent a bad message.\n", waiterID)
			break
		}
		// Log
		x, y, color := parseEvent(p)
		log.Printf("[PLACE] User %s placed pixel at (%d, %d) with color %v\n", waiterID, x, y, color)
		// User has to wait 60 seconds before setting another pixel
		waiter[waiterID] = time.Now().Add(time.Second * time.Duration(env.Timeout))
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(env.Timeout))
		ch <- b
	}
	server.Pl.Close <- i
}

func writeLoop(conn *websocket.Conn, ch chan []byte) {
	for {
		if p, ok := <-ch; ok {
			conn.WriteMessage(websocket.BinaryMessage, p)
		} else {
			break
		}
	}
	conn.Close()
}

func messageHandler(p []byte) error {
	if !server.Pl.SetPixel(parseEvent(p)) {
		return errors.New("invalid placement")
	}
	server.Pl.Msgs <- p
	return nil
}

func parseEvent(b []byte) (int, int, color.Color) {
	if len(b) != 11 {
		return -1, -1, nil
	}
	x := int(binary.BigEndian.Uint32(b))
	y := int(binary.BigEndian.Uint32(b[4:]))
	return x, y, color.NRGBA{b[8], b[9], b[10], 0xFF}
}
