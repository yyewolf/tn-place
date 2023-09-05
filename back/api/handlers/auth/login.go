package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/yyewolf/tn-place/back/internal/env"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/oauth2"
	ogoogle "golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
)

var scopes = []string{
	"https://www.googleapis.com/auth/userinfo.profile",
	"https://www.googleapis.com/auth/userinfo.email",
	// "https://www.googleapis.com/auth/admin.directory.user.readonly",
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
	token := env.C.Cookies.Secret
	key := []byte(token)
	maxAge := 86400 * 30 // 30 days cookie

	store := sessions.NewCookieStore(key)
	store.MaxAge(maxAge)
	store.Options.Domain = env.C.Cookies.Host
	store.Options.Path = "/"

	b, err := os.ReadFile(env.C.Google.Secret)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	s, err := ogoogle.ConfigFromJSON(b, people.ContactsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	googleProvider = google.New(s.ClientID, s.ClientSecret, env.C.Google.RedirectURI, scopes...)
	googleProvider.SetHostedDomain("telecomnancy.net")
	googleProvider.SetPrompt("consent")
	googleProvider.SetAccessType("offline")
	goth.UseProviders(
		googleProvider,
	)
}

func login(ctx *gin.Context) {
	_, err := GetUser(ctx)
	if err != nil {
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		return
	}

	ctx.Redirect(http.StatusFound, env.C.Google.RedirectFront)
}

func GetUser(ctx *gin.Context) (user goth.User, err error) {
	providerI, _ := ctx.Get("provider")
	provider := providerI.(goth.Provider)

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
		s := sess.(*google.Session)
		var token *oauth2.Token

		token, err = provider.RefreshToken(s.RefreshToken)
		if err == nil {
			s.AccessToken = token.AccessToken
			s.RefreshToken = token.RefreshToken
			s.ExpiresAt = token.Expiry
			gothic.StoreInSession(provider.Name(), sess.Marshal(), ctx.Request, ctx.Writer)
		}

		user, err = provider.FetchUser(sess)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"logged": false})
		}
	}

	return
}
