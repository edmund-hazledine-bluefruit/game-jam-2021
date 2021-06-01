package main

import (
	_ "embed"
)

//go:embed static/state_default.json
var gameStateJson []byte

var cards = map[int]Card{
	1: Card{
		Id:          1,
		Name:        "Orange",
		Image:       "orange.jpg",
		Description: "But it's not orange! :o (Special: Draw card - if card is blue, play it)",
		Blue:        true,
		Special:     true,
		Effects: Effects{
			Attack:  0,
			Actions: 1,
			Cards:   1,
		},
	},
	2: Card{
		Id:          2,
		Name:        "Apple",
		Image:       "apple.jpg",
		Description: "It's just an apple.",
		Blue:        false,
		Special:     false,
		Effects: Effects{
			Attack:  1,
			Actions: 1,
			Cards:   1,
		},
	},
	3: Card{
		Id:          3,
		Name:        "Banana",
		Image:       "banana.jpg",
		Description: "This is not a gun.",
		Blue:        true,
		Special:     false,
		Effects: Effects{
			Attack:  2,
			Actions: 0,
			Cards:   1,
		},
	},
	4: Card{
		Id:          4,
		Name:        "Melon",
		Image:       "melon.jpg",
		Description: "Throw this for maximum effect.",
		Blue:        true,
		Special:     false,
		Effects: Effects{
			Attack:  3,
			Actions: 0,
			Cards:   0,
		},
	},
	5: Card{
		Id:          5,
		Name:        "Lime",
		Image:       "lime.jpg",
		Description: "An evil sour lime. Make someone else bite it. (Special: Play alongisde a blue card - Adds +1 (or +2 randomly?) to a random effect)",
		Blue:        false,
		Special:     true,
		Effects: Effects{
			Attack:  0,
			Actions: 0,
			Cards:   0,
		},
	},
	6: Card{
		Id:          6,
		Name:        "Rotten Fruit",
		Image:       "rotten.jpg",
		Description: "Ugh, stinky.",
		Blue:        false,
		Special:     false,
		Effects: Effects{
			Attack:  0,
			Actions: 0,
			Cards:   0,
		},
	},
}

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
	Hand     []Card `json:"hand"`
	Deck     []Card `json:"deck"`
	Discard  []Card `json:"discard"`
	PlayArea []Card `json:"playArea"`
	Score    int    `json:"score"`
	Actions  int    `json:"actions"`
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
