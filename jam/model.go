package main

import (
	_ "embed"
)

//go:embed static/state_default.json
var gameStateJson []byte

type PlayerStateView struct {
	PlayerTurn bool         `json:"playerTurn"`
	Player     PlayerInfo   `json:"player"`
	Opponent   OpponentInfo `json:"opponent"`
	BuyArea    []SupplyPile `json:"buyArea"`
	LastPlayed Card         `json:"lastPlayed"`
}

type PlayerInfo struct {
	Hand    []Card `json:"hand"`
	Deck    int    `json:"deck"`
	Discard int    `json:"discard"`
	Score   int    `json:"score"`
	Actions int    `json:"actions"`
}

type OpponentInfo struct {
	Hand    int `json:"hand"`
	Deck    int `json:"deck"`
	Discard int `json:"discard"`
	Score   int `json:"score"`
}

type GameState struct {
	PlayerOneTurn bool         `json:"playerOneTurn"`
	PlayerOne     Player       `json:"playerOne"`
	PlayerTwo     Player       `json:"playerTwo"`
	BuyArea       []SupplyPile `json:"buyArea"`
	LastPlayed    Card         `json:"lastPlayed"`
}

type Player struct {
	Hand    []Card `json:"hand"`
	Deck    []Card `json:"deck"`
	Discard []Card `json:"discard"`
	Score   int    `json:"score"`
	Actions int    `json:"actions"`
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

func (g *GameState) getPlayerInfo(player string) PlayerStateView {
	p := g.PlayerOne
	o := g.PlayerTwo
	turn := g.PlayerOneTurn

	if player == "2" {
		p = g.PlayerTwo
		o = g.PlayerOne
		turn = !turn
	}

	playerInfo := PlayerInfo{
		Hand:    p.Hand,
		Deck:    len(p.Deck),
		Discard: len(p.Discard),
		Score:   p.Score,
		Actions: p.Actions,
	}

	opponentInfo := OpponentInfo{
		Hand:    len(o.Hand),
		Deck:    len(o.Deck),
		Discard: len(o.Discard),
		Score:   o.Score,
	}

	return PlayerStateView{
		PlayerTurn: turn,
		Player:     playerInfo,
		Opponent:   opponentInfo,
		BuyArea:    g.BuyArea,
		LastPlayed: g.LastPlayed,
	}
}
