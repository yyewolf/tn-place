package auth

import (
	"net/http"

	"github.com/yyewolf/tn-place/back/internal/env"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func callbackRetry(ctx *gin.Context, provider goth.Provider, sess goth.Session) (err error) {
	params := ctx.Request.URL.Query()
	if params.Encode() == "" && ctx.Request.Method == "POST" {
		ctx.Request.ParseForm()
		params = ctx.Request.Form
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, params)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	err = gothic.StoreInSession(provider.Name(), sess.Marshal(), ctx.Request, ctx.Writer)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	_, err = provider.FetchUser(sess)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}

	return
}

func callback(ctx *gin.Context) {
	providerI, _ := ctx.Get("provider")
	provider := providerI.(goth.Provider)

	value, err := gothic.GetFromSession(provider.Name(), ctx.Request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	usr, err := provider.FetchUser(sess)
	if err != nil {
		err = callbackRetry(ctx, provider, sess)
		if err != nil {
			return
		}
	}

	// Disconnect if email doesn't end with @telecomnancy.net
	if usr.Email != "" && usr.Email[len(usr.Email)-16:] != "@telecomnancy.net" {
		gothic.Logout(ctx.Writer, ctx.Request)
		return
	}

	ctx.Redirect(http.StatusFound, env.GoogleRedirectFront)
}
