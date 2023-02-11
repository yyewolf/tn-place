package internal

import (
	"io/ioutil"
	"tn-place/internal/env"
	"tn-place/internal/server"

	"github.com/gin-gonic/gin"
)

func save(c *gin.Context) {
	ioutil.WriteFile(env.SavePath, server.Pl.GetImageBytes(), 0644)
}
