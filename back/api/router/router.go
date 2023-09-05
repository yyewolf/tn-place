package router

import (
	"os"
	"time"

	"github.com/yyewolf/tn-place/back/api/handlers/admin"
	"github.com/yyewolf/tn-place/back/api/handlers/auth"
	"github.com/yyewolf/tn-place/back/api/handlers/gateway"
	"github.com/yyewolf/tn-place/back/api/handlers/image"
	middlewares "github.com/yyewolf/tn-place/back/api/handlers/middleware"
	"github.com/yyewolf/tn-place/back/api/handlers/pixel"
	"github.com/yyewolf/tn-place/back/api/handlers/status"
	"github.com/yyewolf/tn-place/back/internal/canva"
	"github.com/yyewolf/tn-place/back/internal/env"
	"github.com/yyewolf/tn-place/back/internal/server"

	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
)

func Route(engine *gin.Engine) {
	engine.Use(middlewares.WithProvider())
	engine.Use(static.Serve("/", static.LocalFile("dist", false)))
	// engine.NoRoute(static.Serve("/", static.LocalFile("dist", false)))

	path := engine.Group("/")

	// Create image
	cv := canva.NewImage()

	pl := server.NewServer(cv, env.C.ConnectionCount)
	// Watch dog for saving image
	defer os.WriteFile(env.C.SavePath, pl.GetImageBytes(), 0644)
	go func() {
		for {
			os.WriteFile(env.C.SavePath, pl.GetImageBytes(), 0644)
			canva.SavePlacers(pl.Canva.Placers)
			time.Sleep(time.Second * time.Duration(env.C.SaveInterval))
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
