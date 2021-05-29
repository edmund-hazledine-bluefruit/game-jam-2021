package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error
	gameStateJson, err := ioutil.ReadFile("static/state_default.json")
	var gameState GameState
	json.Unmarshal([]byte(gameStateJson), &gameState)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	fmt.Println("Starting Server...")

	router := gin.Default()
	router.StaticFile("/", "./static/home.html")
	router.Static("/static", "./static/")
	router.LoadHTMLGlob("templates/*")

	router.GET("/welcome", handleWelcome)
	router.GET("/wait-for-table", waitForTable)
	router.GET("/game/:game-id", handleGame)

	router.Run()
}
