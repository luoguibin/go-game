package game

import (
	"fmt"
	"go-game/models"
	"net/http"
	"sync"

	"github.com/astaxie/beego/context"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	MGameServer = &GameServer{}

	GameServerStatus = 1
)

type GameServer struct {
	clientMap *sync.Map
}

/*
 * init MGameServer
 */
func init() {
	MGameServer.Start()
}

// Start ...
func (gameServer *GameServer) Start() {
	fmt.Println("GameServer::Start()")
	gameServer.clientMap = &sync.Map{}
}

// AddToServer ...
func AddToServer(Ctx *context.Context, ID int64) {
	ws, err := upgrader.Upgrade(Ctx.ResponseWriter, Ctx.Request, nil)
	if err != nil {
		models.MConfig.MLogger.Error("get ws error:\n%s", err)
	}
	models.MConfig.MLogger.Debug("get ws: " + ws.RemoteAddr().String())

	gameClient := &GameClient{
		ID:   ID,
		Conn: ws,
	}
	data, err := models.QueryGameData(ID)
	if err != nil {
		models.MConfig.MLogger.Error("get ws error:\n%s", err)
		return
	}
	gameClient.Data = data
	MGameServer.clientMap.Store(ID, gameClient)
	go GoGameClientHandle(gameClient)
}

// GoGameClientHandle ...
func GoGameClientHandle(gameClient *GameClient) {
	for {
		// 获取指令
		var order models.GameOrder
		err := gameClient.Conn.ReadJSON(&order)
		if err != nil {
			models.MConfig.MLogger.Error("ws read msg error: " + err.Error())
			return
		}

		order.Data = order.Data.(string) + "(dealed)"

		MGameServer.clientMap.Range(func(key, client_ interface{}) bool {
			client, ok := (client_).(*GameClient)
			if !ok {
				models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
				return true
			}
			client.Conn.WriteJSON(order)
			return true
		})
	}
}
