package router

import (
	"os"
	"time"
	"tn-place/api/handlers/admin"
	"tn-place/api/handlers/auth"
	"tn-place/api/handlers/gateway"
	"tn-place/api/handlers/image"
	"tn-place/api/handlers/pixel"
	"tn-place/api/handlers/status"
	"tn-place/internal/canva"
	"tn-place/internal/env"
	"tn-place/internal/server"

	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
)

func Route(engine *gin.Engine) {
	path := engine.Group("/")

	// Create image
	cv := canva.NewImage()

	pl := server.NewServer(cv, env.ConnectionCount)
	// Watch dog for saving image
	defer os.WriteFile(env.SavePath, pl.GetImageBytes(), 0644)
	go func() {
		for {
			os.WriteFile(env.SavePath, pl.GetImageBytes(), 0644)
			canva.SavePlacers(pl.Canva.Placers)
			time.Sleep(time.Second * time.Duration(env.SaveInterval))
		}
	}()

	server.Pl = pl

	engine.Use(static.Serve("/", static.LocalFile("dist", false)))
	// engine.NoRoute(static.Serve("/", static.LocalFile("dist", false)))

	status.LoadRoutes(path)
	image.LoadRoutes(path)
	gateway.LoadRoutes(path)
	auth.LoadRoutes(path)
	admin.LoadRoutes(path)
	pixel.LoadRoutes(path)
}
