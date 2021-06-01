package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GameEndedMsg struct {
	GameEnded      bool `json:"gameEnded"`
	PlayerOneScore int  `json:"playerOneScore"`
	PlayerTwoScore int  `json:"playerTwoScore"`
}

type Action struct {
	ActionType int `json:"actionType"`
	CardId     int `json:"cardId"`
	UsedWith   int `json:"usedWith"`
}

type Game struct {
	Id            uuid.UUID
	GameState     GameState
	PlayerOneConn *websocket.Conn
	PlayerTwoConn *websocket.Conn
	Ended         bool
}

const (
	OrangeCardId = 1
	LimeCardId   = 5
)

const (
	UseCard = iota
	BuyCard
	EndActions
)

var games map[uuid.UUID]*Game = make(map[uuid.UUID]*Game)

func (game *Game) processAction(action Action, playerId string) {
	var player *Player
	var opponent *Player

	if playerId == "1" {
		player = &game.GameState.PlayerOne
		opponent = &game.GameState.PlayerTwo
	} else {
		player = &game.GameState.PlayerTwo
		opponent = &game.GameState.PlayerOne
	}

	switch action.ActionType {
	case UseCard:
		player.playActionCard(action.CardId, opponent, &game.GameState, action)
	case EndActions:
		player.Actions = 0
	case BuyCard:
		player.buyCard(action.CardId, &game.GameState)
		player.discardAndRedraw()
		if game.gameEnded() {
			return
		}
		game.GameState.PlayerOneTurn = !game.GameState.PlayerOneTurn
		player.Actions = 1
	}

	playerOneInfo := game.GameState.getPlayerInfo("1")
	playerTwoInfo := game.GameState.getPlayerInfo("2")
	playerOneMsg, _ := json.Marshal(playerOneInfo)
	playerTwoMsg, _ := json.Marshal(playerTwoInfo)

	game.PlayerOneConn.WriteMessage(websocket.TextMessage, playerOneMsg)
	game.PlayerTwoConn.WriteMessage(websocket.TextMessage, playerTwoMsg)
	game.GameState.Info = ""
}

func (player *Player) playActionCard(cardId int, opponent *Player, gameState *GameState, action Action) {
	card := cards[cardId]
	gameState.LastPlayed = card
	opponent.Score -= card.Effects.Attack
	player.Actions += card.Effects.Actions
	for i := 0; i < card.Effects.Cards; i++ {
		player.drawCard()
	}

	if card.Special {
		player.handleSpecialCard(card, opponent, gameState, action)
	}

	player.removeCardFromHand(cardId)
	player.Actions--
}

func (player *Player) handleSpecialCard(card Card, opponent *Player, gameState *GameState, action Action) {
	playerName := "Player 1"
	if !gameState.PlayerOneTurn {
		playerName = "Player 2"
	}

	switch card.Id {
	case OrangeCardId:
		player.playOrangeCard(playerName, opponent, gameState)
	case LimeCardId:
		player.playLimeCard(playerName, opponent, gameState, action)
	}
}

func (player *Player) playLimeCard(playerName string, opponent *Player, gameState *GameState, action Action) {
	blueCard := cards[action.UsedWith]
	if !blueCard.Blue {
		gameState.Info = playerName + " played a Lime with a " + blueCard.Name + ", it had no effect."
		return
	}

	rand.Seed(time.Now().UnixNano())
	bonus := rand.Intn(2) + 1
	effect := rand.Intn(3)

	switch effect {
	case 0:
		gameState.Info = playerName + " played a Lime with a " +
			blueCard.Name + ". Adding +" + strconv.Itoa(bonus) + " attack"
		opponent.Score -= bonus
	case 1:
		gameState.Info = playerName + " played a Lime with a " +
			blueCard.Name + ". Adding +" + strconv.Itoa(bonus) + " actions"
		player.Actions += bonus
	case 2:
		gameState.Info = playerName + " played a Lime with a " +
			blueCard.Name + ". Adding +" + strconv.Itoa(bonus) + " card draw"
		for i := 0; i < bonus; i++ {
			player.drawCard()
		}
	}

	player.Actions++
	player.playActionCard(blueCard.Id, player, gameState, action)
}

func (player *Player) playOrangeCard(playerName string, opponent *Player, gameState *GameState) {
	newCard, err := player.drawCard()
	if err != nil {
		gameState.Info = playerName + " played an Orange but there's no cards left to draw"
		return
	}

	if newCard.Blue {
		gameState.Info = playerName + " uses an Orange, they draw and play a " + newCard.Name
		player.playActionCard(newCard.Id, opponent, gameState, Action{})
	} else {
		gameState.Info = playerName + " played an Orange and draws a " + newCard.Name + " better luck next time!"
	}
}

func (player *Player) removeCardFromHand(cardId int) {
	for i, card := range player.Hand {
		if card.Id == cardId {
			player.Hand = append(player.Hand[:i], player.Hand[i+1:]...)
			player.PlayArea = append(player.PlayArea, card)
			return
		}
	}
}

func (player *Player) buyCard(cardId int, gameState *GameState) {
	for i, supplyPile := range gameState.BuyArea {
		if supplyPile.Card.Id != cardId {
			continue
		}

		if supplyPile.Amount > 0 {
			player.Discard = append(player.Discard, supplyPile.Card)
			gameState.BuyArea[i].Amount--
		}

		break
	}
}

func (player *Player) drawCard() (card Card, err error) {
	if len(player.Deck) < 1 {
		player.Deck = shuffle(player.Discard)
		player.Discard = make([]Card, 0, 10)
	}

	if len(player.Deck) < 1 {
		return card, errors.New("Failed to draw card")
	}

	player.Hand = append(player.Hand, player.Deck[:1]...)
	player.Deck = player.Deck[1:]

	return player.Hand[len(player.Hand)-1], err
}

func (player *Player) discardAndRedraw() {
	player.Discard = append(player.Discard, player.PlayArea...)
	player.PlayArea = make([]Card, 0, 10)

	player.Discard = append(player.Discard, player.Hand...)

	if len(player.Deck) >= 5 {
		player.Hand = player.Deck[:5]
		player.Deck = player.Deck[5:]
		return
	}

	player.Hand = player.Deck
	player.Deck = shuffle(player.Discard)
	player.Discard = make([]Card, 0, 10)

	draw := 5 - len(player.Hand)
	player.Hand = append(player.Hand, player.Deck[:draw]...)
	player.Deck = player.Deck[draw:]
}

func (game *Game) gameEnded() (ended bool) {
	count := 0
	for _, supplyPile := range game.GameState.BuyArea {
		if supplyPile.Amount < 1 {
			count++
		}
	}

	if count >= 3 {
		msg := GameEndedMsg{
			GameEnded:      true,
			PlayerOneScore: game.GameState.PlayerOne.Score,
			PlayerTwoScore: game.GameState.PlayerTwo.Score,
		}
		jsonMsg, _ := json.Marshal(msg)
		closeConnections(game, jsonMsg)
		return true
	}
	return false
}

func shuffle(cards []Card) (shuffledCards []Card) {
	if len(cards) < 1 {
		return cards
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	return cards
}

func getGame(gameId uuid.UUID) (game *Game) {
	currentGame, gotGame := games[gameId]
	if gotGame {
		return currentGame
	}

	var gameState GameState
	json.Unmarshal(gameStateJson, &gameState)
	newGame := Game{
		Id:        gameId,
		GameState: gameState,
	}
	games[newGame.Id] = &newGame

	return &newGame
}

func gameSock(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection")
		return
	}
	game, player, err := setupGameConnection(conn, c)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":true}`))
		return
	}

	playerInfo := game.GameState.getPlayerInfo(player)
	out, _ := json.Marshal(playerInfo)

	conn.WriteMessage(websocket.TextMessage, out)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil && !game.Ended {
			closeConnections(game, []byte(`{"error":true}`))
			break
		}
		var action Action
		if json.Unmarshal(msg, &action) == nil {
			game.processAction(action, player)
		}
	}
}

func setupGameConnection(conn *websocket.Conn, c *gin.Context) (game *Game, player string, err error) {
	player = c.Request.URL.Query().Get("player")
	rawGameId := c.Request.URL.Query().Get("gameId")
	gameId, err := uuid.Parse(rawGameId)

	if (player != "1" && player != "2") || err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":true}`))
		conn.Close()
		return
	}

	game = getGame(gameId)

	if player == "1" {
		game.PlayerOneConn = conn
	} else {
		game.PlayerTwoConn = conn
	}

	table, gotTable := tables[gameId]
	if gotTable && player == "2" {
		table.Conn.WriteMessage(websocket.TextMessage, []byte(`{"gameId":"`+rawGameId+`"}`))
	}

	return game, player, err
}

func closeConnections(game *Game, msg []byte) {
	game.Ended = true

	if game.PlayerOneConn != nil {
		game.PlayerOneConn.WriteMessage(websocket.TextMessage, msg)
	}

	if game.PlayerTwoConn != nil {
		game.PlayerTwoConn.WriteMessage(websocket.TextMessage, msg)
	}

	if game.PlayerOneConn != nil {
		game.PlayerOneConn.Close()
		game.PlayerOneConn = nil
	}

	if game.PlayerTwoConn != nil {
		game.PlayerTwoConn.Close()
		game.PlayerTwoConn = nil
	}

	delete(games, game.Id)
}
