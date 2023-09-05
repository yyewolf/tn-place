package auth

import (
	"github.com/gin-gonic/gin"
)

func LoadRoutes(rg *gin.RouterGroup) {
	sg := rg.Group("/auth")
	sg.GET("/", login)
	sg.GET("/self", check)
	sg.GET("/callback", callback)

	lg := rg.Group("/logout")
	lg.GET("/", logout)
}
