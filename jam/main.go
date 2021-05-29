package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting Server...")

	router := gin.Default()
	router.StaticFile("/", "./templates/home.html")
	router.Static("/static", "./static/")
	router.LoadHTMLGlob("templates/*")

	router.GET("/welcome", handleGame)

	router.Run()
}
