package middleware

import (
	"github.com/dhawton/log4g"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/vzau/thoth/internal/server/response"
	"github.com/vzau/thoth/pkg/database"
	dbTypes "github.com/vzau/types/database"
	"gorm.io/gorm/clause"
)

var log = log4g.Category("middleware/auth")

func Auth(c *gin.Context) {
	session := sessions.Default(c)
	cid := session.Get("cid").(uint)

	if cid == 0 {
		c.Set("x-cid", 0)
		c.Set("x-user", nil)
		c.Next()
		return
	}

	user := &dbTypes.User{}
	if err := database.DB.Where(&dbTypes.User{CID: cid}).Preload(clause.Associations).First(&user).Error; err != nil {
		log.Warning("User not found: %d", cid)
		c.Set("x-cid", 0)
		c.Set("x-user", &user)
		response.RespondError(c, 401, "Unauthorized")
		c.Abort()
		return
	}

	c.Set("x-cid", user.CID)
	c.Set("x-user", user)
	c.Next()
}

func NotGuest(c *gin.Context) {
	if c.GetUint("x-cid") == 0 {
		response.RespondError(c, 401, "Unauthorized")
		c.Abort()
		return
	}

	c.Next()
}

func HasRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("x-user").(*dbTypes.User)

		if !userHasRole(user, role) {
			response.RespondError(c, 403, "Forbidden")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func userHasRole(user *dbTypes.User, role string) bool {
	for _, v := range user.Roles {
		if v.Name == role {
			return true
		}
	}

	return false
}

func IsStaff(c *gin.Context) bool {
	user := c.MustGet("x-user").(*dbTypes.User)
	staffRoles := []string{"ATM", "DATM", "TA", "EC", "FE", "WM"}

	for _, v := range staffRoles {
		if userHasRole(user, v) {
			return true
		}
	}

	return false
}
