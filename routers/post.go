package routers

import (
	"blog/apis"
	"blog/middlewares"

	"github.com/gin-gonic/gin"
)

// PostRouterInit - Initialise Post Router
func PostRouterInit(router *gin.Engine) {
	postRouter := router.Group("/post")
	{
		postRouter.POST("/", middlewares.LoginAuth(), apis.AddPost)
		postRouter.PUT("/detail/:pid", middlewares.LoginAuth(), apis.UpdatePostByID)
		postRouter.DELETE("/detail/:pid", middlewares.LoginAuth(), apis.DeletePostByID)
		postRouter.GET("/detail/:pid", middlewares.SoftAuth(), apis.FindPostByID)
		postRouter.GET("/", middlewares.SoftAuth(), apis.FindAllPosts)
	}
}
