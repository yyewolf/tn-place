package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/yyewolf/tn-place/back/api/router"
	"github.com/yyewolf/tn-place/back/internal/env"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	if env.C.LogPath != "" {
		f, err := os.OpenFile(env.C.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	// Set router
	r := gin.Default()

	// Cors allow all
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"X-Internal-Request", "Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.Route(r)

	server := http.Server{
		Addr:    env.C.Port,
		Handler: r,
	}
	log.Fatal(server.ListenAndServe())
}
