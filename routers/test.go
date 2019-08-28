package routers

import (
	"blog/apis"
	"blog/middlewares"

	"github.com/gin-gonic/gin"
)

// TestRouterInit - Router for testing
func TestRouterInit(router *gin.Engine) {
	postRouter := router.Group("/test")
	{
		postRouter.POST("/1", middlewares.SoftAuth(), apis.TestingTheSetAndGetOfContext)

	}
}
