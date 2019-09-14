package game

import (
	"go-game/models"

	"github.com/gorilla/websocket"
)

// GameClient ...
type GameClient struct {
	ID   int64
	Conn *websocket.Conn
	Data *models.GameData
}
