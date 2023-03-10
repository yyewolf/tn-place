package auth

import (
	"net/http"
	"tn-place/internal/env"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func logout(ctx *gin.Context) {
	gothic.Logout(ctx.Writer, ctx.Request)
	http.Redirect(ctx.Writer, ctx.Request, env.GoogleRedirectFront, http.StatusTemporaryRedirect)

	ctx.JSON(http.StatusOK, gin.H{})
}
