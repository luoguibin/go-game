package game

import (
	"fmt"
	"go-game/helper"
	"go-game/models"
	"net/http"
	"strings"
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
	_, ok := MGameServer.clientMap.Load(ID)
	if ok {
		models.MConfig.MLogger.Error("AddToServer repeat ", ID)
		MGameServer.clientMap.Delete(ID)
	}

	ws, err := upgrader.Upgrade(Ctx.ResponseWriter, Ctx.Request, nil)
	if err != nil {
		models.MConfig.MLogger.Error("get ws error:\n%s", err)
	}
	models.MConfig.MLogger.Debug("get ws: " + getCurrentIP(Ctx.Request))

	gameClient := &GameClient{
		ID:   ID,
		Conn: ws,
	}
	data, err := models.QueryGameData(ID)
	if err != nil {
		models.MConfig.MLogger.Error("get ws error:\n%s", err)
		return
	}
	data.OrderMap = make(map[int]models.GameOrder)
	gameClient.Data = data
	MGameServer.clientMap.Store(gameClient.ID, gameClient)

	go GoGameClientHandle(gameClient)
}

func getCurrentIP(r *http.Request) string {
	// 这里也可以通过X-Forwarded-For请求头的第一个值作为用户的ip
	// 但是要注意的是这两个请求头代表的ip都有可能是伪造的
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		// 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}

// GoGameClientHandle ...
func GoGameClientHandle(gameClient *GameClient) {
	// 推送登录角色信息
	MGameServer.writeJSON(models.GameOrder{
		FromGroup: CG_System,
		FromID:    0,
		ToGroup:   CG_Person,
		ToID:      gameClient.ID,
		Type:      CT_Data,
		ID:        CT_Data_Player,
		Data:      gameClient.Data,
	})

	// 推送在线角色信息
	var clientDatas []*models.GameData
	MGameServer.clientMap.Range(func(key, client_ interface{}) bool {
		client, ok := (client_).(*GameClient)
		if !ok {
			models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
			return true
		}
		clientDatas = append(clientDatas, client.Data)
		return true
	})
	MGameServer.writeJSON(models.GameOrder{
		FromGroup: CG_System,
		FromID:    0,
		ToGroup:   CG_Person,
		ToID:      gameClient.ID,
		Type:      CT_Data,
		ID:        CT_Data_Players,
		Data:      clientDatas,
	})

	for {
		// 获取指令
		var order models.GameOrder
		err := gameClient.Conn.ReadJSON(&order)
		if err != nil {
			models.MConfig.MLogger.Error("ws read msg error: " + err.Error())
			return
		}

		// order.Data = order.Data.(string) + "(dealed)"
		order.TimeCreate = helper.GetMillisecond()
		order.TimeCurrent = order.TimeCreate

		if order.ID == CT_Action_Move {
			gameClient.Data.OrderMap[CT_Action_Move] = order
		}

		MGameServer.writeJSON(order)
	}
}

// writeJSON ...
func (server *GameServer) writeJSON(order models.GameOrder) {
	server.clientMap.Range(func(key, client_ interface{}) bool {
		client, ok := (client_).(*GameClient)
		if !ok {
			models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
			return true
		}
		client.Conn.WriteJSON(order)
		return true
	})
}
