package gateway

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"time"

	"github.com/yyewolf/tn-place/back/internal/canva"
	"github.com/yyewolf/tn-place/back/internal/education"
	"github.com/yyewolf/tn-place/back/internal/env"
	"github.com/yyewolf/tn-place/back/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/markbates/goth"
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
		http.Error(c.Writer, "Server full.", http.StatusServiceUnavailable)
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

	gUser := c.MustGet("user").(*goth.User)

	// Send timer to client
	w, ok := waiter[gUser.UserID]
	if !ok {
		return
	}
	s := time.Until(w).Seconds()
	if s < 0 {
		s = 0
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(s))
	ch <- b
}

var waiter = make(map[string]time.Time)

func readLoop(conn *websocket.Conn, i int, c *gin.Context, ch chan []byte) {
	gUser := c.MustGet("user").(*goth.User)
	edu := c.MustGet("education").(*education.Education)
	waiterID := gUser.UserID
	limiter := rateLimiter()
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		if edu.Team == "" {
			log.Printf("[ERR] %s was ignored due to no team.\n", waiterID)
			continue
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
		if server.Pl.Paused {
			log.Printf("[ERR] %s was ignored due to pause.\n", waiterID)
			continue
		}
		if messageHandler(p) != nil {
			log.Printf("[ERR] %s sent a bad message.\n", waiterID)
			break
		}
		// Log
		x, y, color := parseEvent(p)
		log.Printf("[PLACE] User %s on Team %d placed pixel at (%d, %d) with color %v\n", gUser.Email, edu.Team, x, y, color)
		server.Pl.Canva.Placers[x][y] = &canva.PlacerInfo{
			Name: gUser.Name,
			Team: edu.Team,
		}
		// User has to wait 60 seconds before setting another pixel
		waiter[waiterID] = time.Now().Add(time.Second * time.Duration(env.C.Timeout))
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(env.C.Timeout))
		ch <- b
	}
	server.Pl.Close <- i
}

func writeLoop(conn *websocket.Conn, ch chan []byte) {
	// Send amount of clients to all clients
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(server.Pl.ClientAmount()))
	conn.WriteMessage(websocket.BinaryMessage, b)
	server.Pl.Msgs <- b

	for {
		if p, ok := <-ch; ok {
			conn.WriteMessage(websocket.BinaryMessage, p)
		} else {
			break
		}
	}

	conn.Close()

	b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(server.Pl.ClientAmount()))
	server.Pl.Msgs <- b
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
