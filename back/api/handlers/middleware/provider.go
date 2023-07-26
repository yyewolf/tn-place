package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
)

func WithProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provider, err := goth.GetProvider("google")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		values := ctx.Request.URL.Query()
		values.Add("provider", "google")
		ctx.Request.URL.RawQuery = values.Encode()

		ctx.Set("provider", provider)
		ctx.Next()
	}
}
