package admin

import (
	middlewares "github.com/yyewolf/tn-place/back/api/handlers/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.RouterGroup) {
	sg := r.Group("/admin")

	sg.POST("/resize", middlewares.IsInternal(), resize)
	sg.GET("/save", middlewares.IsInternal(), save)
	sg.GET("/pause", pause)
	sg.POST("/pause", middlewares.IsInternal(), pause)
}
