package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Port string
var LoadPath string
var SavePath string
var SaveInterval int
var LogPath string
var Width int
var Height int
var ConnectionCount int
var Timeout int
var GoogleSecretPath string
var GoogleRedirectURI string
var GoogleRedirectFront string
var CookieHost string
var CookieSecret string
var InternalSecret string

func init() {
	godotenv.Load()
	var err error
	Port = os.Getenv("PORT")
	LoadPath = os.Getenv("LOAD")
	SavePath = os.Getenv("SAVE")
	LogPath = os.Getenv("LOG")
	GoogleSecretPath = os.Getenv("GOOGLE_SECRET")
	GoogleRedirectURI = os.Getenv("GOOGLE_REDIRECT_URI")
	GoogleRedirectFront = os.Getenv("GOOGLE_REDIRECT_FRONT")
	CookieHost = os.Getenv("COOKIE_HOST")
	CookieSecret = os.Getenv("COOKIE_SECRET")
	InternalSecret = os.Getenv("INTERNAL_SECRET")
	Width, err = strconv.Atoi(os.Getenv("WIDTH"))
	if err != nil {
		log.Fatal(err)
	}
	Height, err = strconv.Atoi(os.Getenv("HEIGHT"))
	if err != nil {
		log.Fatal(err)
	}
	ConnectionCount, err = strconv.Atoi(os.Getenv("COUNT"))
	if err != nil {
		log.Fatal(err)
	}
	SaveInterval, err = strconv.Atoi(os.Getenv("SAVE_INTERVAL"))
	if err != nil {
		log.Fatal(err)
	}
	Timeout, err = strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}
}
