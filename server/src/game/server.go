package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	Wait        sync.WaitGroup
	BanWordBase []string
	Lock        *sync.RWMutex

	ServerConfig *serverConfig
}

type serverConfig struct {
	Serverid int64     `json:"serverid"`
	Host     int       `json:"host"`
	Savepath string    `json:"savepath"`
	Database *Database `json:"database"`
}

type Database struct {
	Databasetype string `json:"databasetype"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

var server *Server
var HttpServer http.Server

func GetServer() *Server {
	if server == nil {
		server = new(Server)
		server.Lock = new(sync.RWMutex)
		server.ServerConfig = new(serverConfig)
		server.ServerConfig.Database = new(Database)
	}
	return server
}

var playerGM *Player

func (self *Server) Start() {
	rand.Seed(time.Now().Unix()) //开启一个随机数种子
	self.LoadConfig()            //加载服务器配置文件
	csvs.CheckCsvs()             //加载游戏数据
	fmt.Println("游戏测试开始")
	go GetManagerBanWord().Run()            //开启一个协程，违禁词汇
	playerGM = NewTestPlayer(nil, 20221019) //创建一个玩家，通常是客户端发来
	go playerGM.Run()                       //玩家协程开始运行
	self.Wait.Wait()
	// go GetManagerHttp().InitData()//开启http协议，并且初始化数据
	// Initservlet(self.ServerConfig.Host)//开启端口
}

func (self *Server) Close() {
	GetManagerBanWord().Close()//关闭违禁词汇协程
	HttpServer.Close()//关闭http服务
	fmt.Println("服务器关闭成功")
}

func (self *Server) GoAdd() {
	self.Wait.Add(1)
}
func (self *Server) GoDone() {
	self.Wait.Done()
}
func Initservlet(port int) {
	HttpServer = http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: nil,
	}
	fmt.Printf("HTTP SERVER STARTING AND LISTENING AT: %d\n", port)
	err := HttpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("HTTP SERVER CLOSE %v\n", err.Error())
	}

}
func (self *Server) IsBanWord(txt string) bool {
	self.Lock.RLock()
	defer self.Lock.Unlock()
	for _, v := range self.BanWordBase {
		matche, err := regexp.MatchString(v, txt)
		if err != nil {
			fmt.Println("IsBanWord err:" + err.Error())
			return false
		}
		if matche {
			return matche
		}
	}
	return false
}
func (self *Server) UpBanWords(txt []string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	self.BanWordBase = txt
}

func (self *Server) Signal() {
	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, syscall.SIGINT)
	for {
		select {
		//放入参数，重写了这个参数代表的操作
		case <-channelSignal:
			self.Close()
		}
	}
}

func (self *Server) LoadConfig() {
	config, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("配置文件读取错误")
		return
	}
	err = json.Unmarshal(config, self.ServerConfig)
	if err != nil {
		return
	}
}
