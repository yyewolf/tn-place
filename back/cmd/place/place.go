package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yyewolf/tn-place/back/api/router"
	"github.com/yyewolf/tn-place/back/internal/env"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	if env.LogPath != "" {
		f, err := os.OpenFile(env.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	// Set router
	r := gin.Default()

	// Cors allow all
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	})

	router.Route(r)

	server := http.Server{
		Addr:    env.Port,
		Handler: r,
	}
	log.Fatal(server.ListenAndServe())
}
