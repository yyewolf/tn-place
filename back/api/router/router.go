package router

import (
	"io/ioutil"
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

	"github.com/gin-gonic/gin"
)

func Route(engine *gin.Engine) {
	path := engine.Group("/")

	// Create image
	cv := canva.NewImage()

	pl := server.NewServer(cv, env.ConnectionCount)
	// Watch dog for saving image
	defer ioutil.WriteFile(env.SavePath, pl.GetImageBytes(), 0644)
	go func() {
		for {
			ioutil.WriteFile(env.SavePath, pl.GetImageBytes(), 0644)
			canva.SavePlacers(pl.Canva.Placers)
			time.Sleep(time.Second * time.Duration(env.SaveInterval))
		}
	}()
	server.Pl = pl

	status.LoadRoutes(path)
	image.LoadRoutes(path)
	gateway.LoadRoutes(path)
	auth.LoadRoutes(path)
	admin.LoadRoutes(path)
	pixel.LoadRoutes(path)
}
