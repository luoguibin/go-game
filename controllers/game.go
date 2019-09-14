package controllers

import (
	"go-game/game"
	"strconv"
)

// GameController ...
type GameController struct {
	BaseController
}

// Get WebSocket连接入口，在BeforeRouter检测jwt中的合法后才给予长连接
func (c *GameController) Get() {
	ID, _ := strconv.ParseInt(c.Ctx.Input.Query("uId"), 10, 64)
	game.AddToServer(c.Ctx, ID)
}
