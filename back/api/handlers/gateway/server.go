package gateway

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"time"

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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &server.PlaceClient{
		LastHeartbeat: time.Now(),
		Channel:       make(chan []byte, 8),
	}
	server.Pl.Clients = append(server.Pl.Clients, client)

	loggedInInterface, _ := c.Get("is_logged_in")
	loggedIn := loggedInInterface.(bool)
	if !loggedIn {
		conn.WriteMessage(websocket.BinaryMessage, []byte("not_logged_in"))
	}

	go readLoop(conn, c, client, loggedIn)
	go writeLoop(conn, client)

	if !loggedIn {
		return
	}

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
	client.Channel <- b
}

var waiter = make(map[string]time.Time)

func readLoop(conn *websocket.Conn, c *gin.Context, client *server.PlaceClient, loggedIn bool) {
	var gUser *goth.User
	var waiterID string
	if loggedIn {
		gUser = c.MustGet("user").(*goth.User)
		waiterID = gUser.UserID
	}

	limiter := rateLimiter()
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		// Check for heartbeat
		if string(p) == "hb" {
			client.LastHeartbeat = time.Now()
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
		log.Printf("[PLACE] User %s placed pixel at (%d, %d) with color %v\n", waiterID, x, y, color)

		// server.Pl.Canva.Placers[x][y] = gUser.Email
		server.Pl.Canva.Placers[x][y] = "Anonymous"

		// User has to wait 60 seconds before setting another pixel
		waiter[waiterID] = time.Now().Add(time.Second * time.Duration(env.C.Timeout))
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(env.C.Timeout))

		if client.Closed {
			break
		}

		client.Channel <- b
	}

	if !client.Closed {
		server.Pl.Close <- client
	}
}

func writeLoop(conn *websocket.Conn, client *server.PlaceClient) {
	// Send amount of clients to all clients
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(server.Pl.ClientAmount()))
	conn.WriteMessage(websocket.BinaryMessage, b)
	server.Pl.Msgs <- b

	for {
		if client.Closed {
			break
		}

		if p, ok := <-client.Channel; ok {
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
