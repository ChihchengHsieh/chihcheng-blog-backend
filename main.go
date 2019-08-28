package main

import (
	"blog/databases"
	"blog/routers"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	gin.ForceConsoleColor()
	databases.InitDB()
	router := routers.InitRouter()
	router.Run()
}
