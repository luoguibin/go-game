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
	"github.com/goinggo/mapstructure"
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
	client, ok := MGameServer.clientMap.Load(ID)
	if ok {
		models.MConfig.MLogger.Error("AddToServer repeat ", ID)
		MGameServer.removePlayer((client).(*GameClient))
		Ctx.WriteString("已在线，请稍后登录")
		return
	}

	ws, err := upgrader.Upgrade(Ctx.ResponseWriter, Ctx.Request, nil)
	if err != nil {
		models.MConfig.MLogger.Error("get ws error:\n%s", err)
	}
	models.MConfig.MLogger.Debug("get ws: %s", getCurrentIP(Ctx.Request))

	gameClient := &GameClient{
		ID:   ID,
		Conn: ws,
	}
	data, err := models.QueryGameData(ID)
	if err != nil {
		models.MConfig.MLogger.Error("get ws error:\n%s", err)
		return
	}
	data.OrderMap = make(map[int]*models.GameOrder)
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
	order := models.GameOrder{
		FromGroup:  CG_System,
		FromID:     0,
		ToGroup:    CG_Person,
		ToID:       gameClient.ID,
		Type:       CT_Data,
		ID:         CT_Data_Player,
		Data:       gameClient.Data,
		TimeCreate: helper.GetMillisecond(),
	}
	order.TimeCurrent = order.TimeCreate
	MGameServer.braodcastOrder(order)

	// 推送在线角色信息，不包括自身
	var clientDatas []*models.GameData
	MGameServer.clientMap.Range(func(key, client_ interface{}) bool {
		client, ok := (client_).(*GameClient)
		if !ok {
			models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
			return true
		}

		if client.ID == gameClient.ID {
			return true
		}

		time := helper.GetMillisecond()
		for _, order := range client.Data.OrderMap {
			order.TimeCurrent = time
		}

		clientDatas = append(clientDatas, client.Data)
		return true
	})
	order.ID = CT_Data_Players
	order.Data = clientDatas
	gameClient.Conn.WriteJSON(order)

	for {
		// 获取指令
		var order models.GameOrder
		err := gameClient.Conn.ReadJSON(&order)
		if err != nil {
			models.MConfig.MLogger.Error("ws read msg error: " + err.Error())
			MGameServer.removePlayer(gameClient)
			return
		}

		order.TimeCreate = helper.GetMillisecond()
		order.TimeCurrent = order.TimeCreate

		if order.ID == CT_Action_Move {
			MGameServer.updatePlayer(gameClient, false)
			gameClient.Data.OrderMap[CT_Action_Move] = &order
		}

		MGameServer.braodcastOrder(order)
	}
}

// braodcastOrder ...
func (gameServer *GameServer) braodcastOrder(order models.GameOrder) {
	gameServer.clientMap.Range(func(key, client_ interface{}) bool {
		client, ok := (client_).(*GameClient)
		if !ok {
			models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
			return true
		}
		client.Conn.WriteJSON(order)
		return true
	})
}

func (gameServer *GameServer) updatePlayer(client *GameClient, isSave bool) {
	for _, order := range client.Data.OrderMap {
		if order.ID == CT_Action_Move {
			var position Position
			err :=  mapstructure.Decode(order.Data, &position)
			if err == nil {
				client.Data.X = position.X
				// client.Data.Y = position.Y
				client.Data.Z = position.Z
			}
		}
	}
	if isSave {
		models.UpdateGameData(client.Data)
	}
}

// removePlayer ...
func (gameServer *GameServer) removePlayer(client *GameClient) {
	gameServer.updatePlayer(client, true)
	gameServer.clientMap.Delete(client.ID)
	gameServer.braodcastOrder(models.GameOrder{
		FromGroup: CG_System,
		FromID:    0,
		ToGroup:   CG_Person,
		ToID:      client.ID,
		Type:      CT_Data,
		ID:        CT_Data_Remove,
	})
}
