package apihandler

import (
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	JsonOk(c, gin.H{})
}
