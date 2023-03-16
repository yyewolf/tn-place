package middlewares

import (
	"github.com/yyewolf/tn-place/back/api/handlers/auth"

	"github.com/gin-gonic/gin"
)

func SetStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gu, err := auth.GetUser(ctx)
		if err != nil {
			ctx.Set("is_logged_in", false)
			return
		}

		// user, err := services.GetUserByID(gu.UserID)
		// if err != nil {
		// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting user"})
		// }

		ctx.Set("is_logged_in", true)
		ctx.Set("user_id", gu.UserID)
		ctx.Set("user", &gu)
	}
}
