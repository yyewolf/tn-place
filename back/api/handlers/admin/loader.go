package admin

import (
	middlewares "tn-place/api/handlers/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.RouterGroup) {
	sg := r.Group("/admin")

	sg.POST("/resize", middlewares.IsInternal(), resize)
	sg.GET("/save", middlewares.IsInternal(), save)
}
