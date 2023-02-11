package internal

import (
	middlewares "tn-place/api/handlers/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.RouterGroup) {
	sg := r.Group("/internal")

	sg.GET("/resize", middlewares.IsInternal(), resize)
}
