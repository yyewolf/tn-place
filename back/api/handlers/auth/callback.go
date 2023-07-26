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

	_, err = provider.FetchUser(sess)
	if err != nil {
		err = callbackRetry(ctx, provider, sess)
		if err != nil {
			return
		}
	}

	// c := context.Background()
	// oconfig := &oauth2.Config{}
	// token, err := googleProvider.RefreshToken(user.RefreshToken)
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error refreshing token"})
	// 	return
	// }
	// adminService, err := admin.NewService(c, option.WithTokenSource(oconfig.TokenSource(c, token)))
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	// 	return
	// }

	// t, err := adminService.Users.Get(user.UserID).Projection("custom").CustomFieldMask("Education").ViewType("domain_public").Do()
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	// 	return
	// }
	// edc := &Education{}
	// err = json.Unmarshal(t.CustomSchemas["Education"], edc)
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	// 	return
	// }
	// Saves to database

	// alreadyRegistered, _ := services.GetUserByID(user.UserID)

	// u := &models.User{
	// 	ID:      user.UserID,
	// 	Email:   user.Email,
	// 	Name:    user.FirstName,
	// 	Surname: user.LastName,
	// 	Promo:   edc.Promo,
	// 	Spe:     edc.Spe,

	// 	Type: alreadyRegistered.Type,
	// }

	// db.DB.Save(u)

	ctx.Redirect(http.StatusFound, env.GoogleRedirectFront)
}
