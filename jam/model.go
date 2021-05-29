package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var tables map[uuid.UUID]Table = make(map[uuid.UUID]Table)
var gameStateJson string

type Game struct {
	Value int
	Tag   string `json:"label"`
}

type Table struct {
	Id   uuid.UUID
	Name string
	Conn *websocket.Conn
}

type GameState struct {
	PlayerTurn bool         `json:"playerTurn"`
	Player     Player       `json:"player"`
	Opponent   Opponent     `json:"opponent"`
	BuyArea    []SupplyPile `json:"buyArea"`
}

type Player struct {
	Hand    []Card `json:"hand"`
	Desk    int    `json:"deck"`
	Discard int    `json:"discard"`
	Score   int    `json:"score"`
}

type Opponent struct {
	Hand    int `json:"hand"`
	Deck    int `json:"deck"`
	Discard int `json:"discard"`
	Score   int `json:"score"`
}

type SupplyPile struct {
	Amount int  `json:"amount"`
	Card   Card `json:"card"`
}

type Card struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Blue        bool    `json:"blue"`
	Special     bool    `json:"special"`
	Effects     Effects `json:"effects"`
}

type Effects struct {
	Attack  int `json:"attack"`
	Actions int `json:"actions"`
	Cards   int `json:"cards"`
}
