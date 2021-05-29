package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleGame(c *gin.Context) {
	gameId, _ := uuid.Parse(c.Param("game-id"))
	player := c.Request.URL.Query().Get("user-id")
	_, gotTable := tables[gameId]

	if !gotTable || (player != "1" && player != "2") {
		fmt.Println("Failed check")
		fmt.Println(gotTable)
		fmt.Println(player)
		c.Redirect(http.StatusMovedPermanently, "/welcome")
		return
	}

	c.File("static/home.html")
}
