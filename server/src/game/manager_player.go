package game

import (
	"sync"

	"golang.org/x/net/websocket"
)

var managerPlayer *ManagerPlayer

type ManagerPlayer struct {
	Players map[int64]*Player
	lock    *sync.RWMutex
}

func GetManagerPlayer() *ManagerPlayer {
	if managerPlayer == nil {
		managerPlayer = new(ManagerPlayer)
		managerPlayer.Players = make(map[int64]*Player)
		managerPlayer.lock = new(sync.RWMutex)
	}
	return managerPlayer
}

func (self *ManagerPlayer) PlayerLoginIn(ws *websocket.Conn, userId int64) *Player {
	self.lock.Lock()
	defer self.lock.Unlock()
	playerInfo, ok := self.Players[userId]
	if ok {
		//顶号操作
		if ws != playerInfo.Ws {
			if playerInfo.Ws != nil {
				playerInfo.Ws.Write([]byte("账号连接断开"))
				playerInfo.Ws.Close()
				playerInfo.Ws = ws
			}
		}
	}
	playerInfo = NewTestPlayer(ws, userId)
	managerPlayer.Players[playerInfo.UserId] = playerInfo
	return playerInfo
}

func (self *ManagerPlayer) PlayerClose(ws *websocket.Conn, userId int64) {
	self.lock.Lock()
	defer self.lock.Unlock()
	playerInfo, ok := self.Players[userId]
	if ok {
		if ws != playerInfo.Ws {
			playerInfo.Ws.Write([]byte("websockt连接断开,ws置空"))
			playerInfo.Ws = nil
		}
	}
}
