package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Table struct {
	Id   uuid.UUID
	Name string
	Conn *websocket.Conn
}

var tables map[uuid.UUID]Table = make(map[uuid.UUID]Table)
