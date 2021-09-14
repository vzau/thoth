package middleware

import (
	"time"

	"github.com/dhawton/log4g"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Make change that updates cookie value (obfuscation), and also updates maxage
func UpdateCookie(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("u", time.Now().String())
	err := session.Save()
	if err != nil {
		log4g.Category("middleware/cookie").Error("Error saving cookie %s", err.Error())
	}
	c.Next()
}
