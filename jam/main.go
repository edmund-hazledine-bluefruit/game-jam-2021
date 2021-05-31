package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting Server...")

	router := gin.Default()
	router.StaticFile("/play", "./static/home.html")
	router.Static("/static", "./static/")
	router.LoadHTMLGlob("templates/*")
	router.GET("/", handleWelcome)
	router.GET("/wait-for-table", waitForTable)
	router.GET("/gamesocket", gameSock)

	router.Run()
}
