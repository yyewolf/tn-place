package gateway

import "github.com/gin-gonic/gin"

func LoadRoutes(r *gin.RouterGroup) {
	gw := r.Group("/gateway")

	gw.GET("/", GetGateway)
}
