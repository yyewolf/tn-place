package gateway

import (
	middlewares "github.com/yyewolf/tn-place/back/api/handlers/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.RouterGroup) {
	sg := r.Group("/gateway")

	sg.GET("/", middlewares.SetStatus(), GetGateway)
}
