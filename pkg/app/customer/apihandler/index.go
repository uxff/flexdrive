package handler

import (
	"os"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	hostName, _ := os.Hostname()
	JsonOk(c, gin.H{
		"hostname": hostName,
		"info":     "this is flexdrive customer api server",
	})
}
