package routers

import (
	"github.com/gin-gonic/gin"
)

// InitRouter - Initialise all the router in this function
func InitRouter() *gin.Engine {
	router := gin.Default()
	PostRouterInit(router)
	UserRouterInit(router)
	TestRouterInit(router)
	return router
}
