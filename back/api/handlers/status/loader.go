package status

import (
	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.RouterGroup) {
	sg := r.Group("/status")

	sg.GET("/", GetStatus)

}
