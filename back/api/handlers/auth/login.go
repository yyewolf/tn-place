package auth

import (
	"log"
	"net/http"
	"os"
	"tn-place/internal/env"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	ogoogle "golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
)

var scopes = []string{
	"https://www.googleapis.com/auth/userinfo.profile",
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/admin.directory.user.readonly",
}

type Education struct {
	Promo  int    `json:"Promotion"`
	Spe    string `json:"Approfondissement"`
	Statut int    `json:"Statut"`
}

var (
	googleProvider *google.Provider
)

func init() {
	token := env.CookieSecret
	key := []byte(token)
	maxAge := 86400 * 30 // 30 days cookie

	store := sessions.NewCookieStore(key)
	store.MaxAge(maxAge)
	store.Options.Domain = env.CookieHost
	store.Options.Path = "/"

	b, err := os.ReadFile(env.GoogleSecretPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	s, err := ogoogle.ConfigFromJSON(b, people.ContactsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	googleProvider = google.New(s.ClientID, s.ClientSecret, env.GoogleRedirectURI, scopes...)
	googleProvider.SetHostedDomain("telecomnancy.net")
	googleProvider.SetPrompt("consent")
	googleProvider.SetAccessType("offline")
	goth.UseProviders(
		googleProvider,
	)
}

func login(ctx *gin.Context) {
	// try to get the user without re-authenticating
	provider, err := goth.GetProvider("google")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	values := ctx.Request.URL.Query()
	values.Add("provider", "google")
	ctx.Request.URL.RawQuery = values.Encode()

	value, err := gothic.GetFromSession(provider.Name(), ctx.Request)
	if err != nil {
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		return
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		return
	}

	_, err = provider.FetchUser(sess)
	if err != nil {
		params := ctx.Request.URL.Query()
		if params.Encode() == "" && ctx.Request.Method == "POST" {
			ctx.Request.ParseForm()
			params = ctx.Request.Form
		}

		// get new token and retry fetch
		_, err = sess.Authorize(provider, params)
		if err != nil {
			gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
			return
		}

		err = gothic.StoreInSession(provider.Name(), sess.Marshal(), ctx.Request, ctx.Writer)

		if err != nil {
			gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
			return
		}

		_, err := provider.FetchUser(sess)
		if err != nil {
			gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
			return
		}
	}

	ctx.Redirect(http.StatusFound, env.GoogleRedirectFront)
}

func GetUser(ctx *gin.Context) (user goth.User, err error) {
	// try to get the user without re-authenticating
	provider, err := goth.GetProvider("google")
	if err != nil {
		return
	}
	values := ctx.Request.URL.Query()
	values.Add("provider", "google")
	ctx.Request.URL.RawQuery = values.Encode()

	value, err := gothic.GetFromSession(provider.Name(), ctx.Request)
	if err != nil {
		return
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return
	}

	user, err = provider.FetchUser(sess)
	if err != nil {
		params := ctx.Request.URL.Query()
		if params.Encode() == "" && ctx.Request.Method == "POST" {
			ctx.Request.ParseForm()
			params = ctx.Request.Form
		}

		// get new token and retry fetch
		_, err = sess.Authorize(provider, params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"logged": false})
			return
		}

		err = gothic.StoreInSession(provider.Name(), sess.Marshal(), ctx.Request, ctx.Writer)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"logged": false})
			return
		}

		user, err = provider.FetchUser(sess)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"logged": false})
			return
		}
	}

	token, err := provider.RefreshToken(user.RefreshToken)
	if err == nil {
		user.AccessToken = token.AccessToken
		user.RefreshToken = token.RefreshToken
		gothic.StoreInSession(provider.Name(), sess.Marshal(), ctx.Request, ctx.Writer)
	}
	return
}
