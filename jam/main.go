package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting Server...")

	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.StaticFile("/", "./static/home.html")
	router.Static("/static", "./static/")
	router.LoadHTMLGlob("templates/*")

	router.GET("/welcome", handleWelcome)

	router.GET("/wait-for-table", waitForTable)

	router.Run()
}
