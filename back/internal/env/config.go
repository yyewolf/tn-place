package env

import (
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	Google struct {
		Secret        string `env:"SECRET"`
		RedirectURI   string `env:"REDIRECT_URI"`
		RedirectFront string `env:"REDIRECT_FRONT"`
	} `envPrefix:"GOOGLE_"`

	Cookies struct {
		Host   string `env:"HOST"`
		Secret string `env:"SECRET"`
	} `envPrefix:"COOKIE_"`

	InternalSecret string `env:"INTERNAL_SECRET"`
	Port           string `env:"PORT"`
	LoadPath       string `env:"LOAD"`
	SavePath       string `env:"SAVE"`
	LogPath        string `env:"LOG"`

	Width           int `env:"WIDTH"`
	Height          int `env:"HEIGHT"`
	ConnectionCount int `env:"COUNT"`
	SaveInterval    int `env:"SAVE_INTERVAL"`
	Timeout         int `env:"TIMEOUT"`
}

var C Config

func GetConfig() Config {
	return C
}

func init() {
	godotenv.Load()
	if err := env.Parse(&C); err != nil {
		log.Fatal(err)
	}
}
