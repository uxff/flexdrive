package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	//c.HTML(http.StatusOK, "pkg/app/admin/view/login/login.tpl", gin.H{})
	c.HTML(http.StatusOK, "index/index.tpl", gin.H{
		"path": "login",
	})
}
