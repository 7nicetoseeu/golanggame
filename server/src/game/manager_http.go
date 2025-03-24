package game

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var managerHttp *ManagerHttp

type ManagerHttp struct {
	ManagerChannel chan int
}

func GetManagerHttp() *ManagerHttp {
	if managerHttp == nil {
		managerHttp = new(ManagerHttp)
		managerHttp.ManagerChannel = make(chan int)
	}
	return managerHttp
}

func (self *ManagerHttp) InitData() {
	http.Handle("/", websocket.Handler(self.WebSocketHandler))
	http.HandleFunc("/cname", self.CorrectName)
}
func (self *ManagerHttp) CorrectName(w http.ResponseWriter, r *http.Request) {
	playerGM.RecvSetName("newname")
}

func (self *ManagerHttp) WebSocketHandler(ws *websocket.Conn) {
	var player *Player
	fmt.Println("服务器连接成功")
	for {
		var msg []byte
		ws.SetReadDeadline(time.Now().Add(time.Second))
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			fmt.Printf("err.Error(): %v\n", err.Error())
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				continue
			}
			if player != nil {
				GetManagerPlayer().PlayerClose(ws, player.UserId)
			}
			fmt.Println("websockt连接断开")
			break
		}
		if player == nil {
			var msgLogin MsgLogin
			err = json.Unmarshal(msg, &msgLogin)
			if err != nil {
				return
			}
			player = GetManagerPlayer().PlayerLoginIn(ws, msgLogin.UserId)
		}
		fmt.Println(string(msg))
	}
}
