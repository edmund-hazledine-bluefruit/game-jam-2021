package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type game struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Table struct {
	Id   uuid.UUID
	Name string
	Conn *websocket.Conn
}

var tables map[uuid.UUID]Table = make(map[uuid.UUID]Table)

var upgrader = websocket.Upgrader{}

func handleWelcome(c *gin.Context) {
	games := make([]game, len(tables))
	for _, t := range tables {
		games = append(games, game{
			Name: t.Name,
			Id:   t.Id.String(),
		})
	}

	c.HTML(
		http.StatusOK,
		"welcome.html",
		gin.H{
			"Games": games,
		},
	)
}

func waitForTable(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	tableName := c.Request.URL.Query().Get("name")
	if tableName == "" {
		fmt.Println("Invalid table name")
		conn.Close()
	}

	table := Table{
		Id:   uuid.New(),
		Name: tableName,
		Conn: conn,
	}

	tables[table.Id] = table

	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			delete(tables, table.Id)
			break
		}
	}
}
