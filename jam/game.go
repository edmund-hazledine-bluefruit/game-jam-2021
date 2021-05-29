package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type game struct {
	Name string `json:"name"`
}

func handleGame(c *gin.Context) {
	games := []game{game{Name: "Foo Bar"}}

	c.HTML(
		http.StatusOK,
		"welcome.html",
		gin.H{
			"Games": games,
		},
	)
}
