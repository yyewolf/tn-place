package middlewares

import (
	"net/http"

	"github.com/yyewolf/tn-place/back/internal/env"

	"github.com/gin-gonic/gin"
)

func IsInternal() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Header.Get("X-Internal-Request") != env.InternalSecret {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.Next()
	}
}
