package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var tables map[uuid.UUID]Table = make(map[uuid.UUID]Table)

type Game struct {
	Value int
	Tag   string `json:"label"`
}

type GameState struct {
}

type Table struct {
	Id   uuid.UUID
	Name string
	Conn *websocket.Conn
}
