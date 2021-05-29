package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Game struct {
	Id   uuid.UUID
	Game GameState
}

var games map[uuid.UUID]Game = make(map[uuid.UUID]Game)

func createGame(c *gin.Context) {
	//gameId, _ := uuid.Parse(c.Param("game-id"))
	//player := c.Request.URL.Query().Get("user-id")
	//_, gotTable := tables[gameId]

	//if !gotTable || (player != "1" && player != "2") {
	//	c.Redirect(http.StatusMovedPermanently, "/welcome")
	//	return
	//}

	c.File("static/home.html")
}

func gameSock(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	var gameState GameState
	json.Unmarshal(gameStateJson, &gameState)
	r := gameState.getPlayerInfo("1")
	out, _ := json.Marshal(r)
	conn.WriteMessage(websocket.TextMessage, out)
	//tableName := c.Request.URL.Query().Get("name")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Connection dropped!!!", err)
			conn.Close()
			//delete(tables, table.Id)
			break
		}
		fmt.Println(msg)
		conn.WriteMessage(websocket.TextMessage, out)
	}
}
