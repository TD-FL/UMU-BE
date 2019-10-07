package api

import (
	"github.com/gin-gonic/gin"
)

var PrivacyText = ""

func GetPrivacy(c *gin.Context) {
	c.String(200, "", PrivacyText)
}
