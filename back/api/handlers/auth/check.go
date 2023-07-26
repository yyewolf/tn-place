package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func check(ctx *gin.Context) {
	providerI, _ := ctx.Get("provider")
	provider := providerI.(goth.Provider)

	value, err := gothic.GetFromSession(provider.Name(), ctx.Request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"logged": false})
		return
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"logged": false})
		return
	}

	user, err := provider.FetchUser(sess)
	if err != nil {
		s := sess.(*google.Session)

		token, err := provider.RefreshToken(s.RefreshToken)
		if err == nil {
			s.AccessToken = token.AccessToken
			s.RefreshToken = token.RefreshToken
			s.ExpiresAt = token.Expiry
			gothic.StoreInSession(provider.Name(), sess.Marshal(), ctx.Request, ctx.Writer)
		}

		user, err = provider.FetchUser(sess)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"logged": false})
			return
		}

		// u, err := services.GetUserByID(user.UserID)
		// if err != nil {
		// 	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"logged": false})
		// 	return
		// }
		ctx.JSON(http.StatusOK, gin.H{"user": user, "logged": true})
		return
	}

	// u, err := services.GetUserByID(user.UserID)
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"logged": false})
	// 	return
	// }
	ctx.JSON(http.StatusOK, gin.H{"user": user, "logged": true})
}
