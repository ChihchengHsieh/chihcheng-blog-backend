package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TestingTheSetAndGetOfContext - testing the c.set and c.get
func TestingTheSetAndGetOfContext(c *gin.Context) {
	user, ok := c.Get("user")

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot retrieve the user in gin context",
			"msg":   "Cannot retrieve the user in gin context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"ok":   ok,
	})
}
