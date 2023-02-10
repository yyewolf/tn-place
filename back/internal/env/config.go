package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Port string
var Root string
var LoadPath string
var SavePath string
var SaveInterval int
var LogPath string
var Width int
var Height int
var ConnectionCount int
var Timeout int

func init() {
	godotenv.Load()
	var err error
	Port = os.Getenv("PORT")
	Root = os.Getenv("ROOT")
	LoadPath = os.Getenv("LOAD")
	SavePath = os.Getenv("SAVE")
	LogPath = os.Getenv("LOG")
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
