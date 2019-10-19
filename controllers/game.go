package controllers

import (
	"fmt"
	"go-game/game"
)

// GameController ...
type GameController struct {
	BaseController
}

// Get WebSocket连接入口，在BeforeRouter检测jwt中的合法后才给予长连接
func (c *GameController) Get() {
	userID := c.Ctx.Input.GetData("userId").(int64)
	userName := c.Ctx.Input.GetData("userName").(string)
	fmt.Println("GameController", userID)
	game.AddToServer(c.Ctx, userID, userName)
}
