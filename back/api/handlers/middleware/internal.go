package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yyewolf/tn-place/back/internal/env"
)

func IsInternal() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Header.Get("X-Internal-Request") != env.C.InternalSecret {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.Next()
	}
}
